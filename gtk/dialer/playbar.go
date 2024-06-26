package dialer

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"google.golang.org/grpc/status"
)

// ApplyPlaylistGroupChanges applies collection changes.
func ApplyPlaylistGroupChanges() {
	slog.Info("Applying playlist-group changes")

	requests, err := store.BData.GetUpdatePlaylistGroupRequests()
	if err != nil {
		slog.Error("Failed to get update-playlist-group requests", "error", err)
	}

	cc, err := getClientConn1()
	if err != nil {
		slog.Error("Failed to get client connection", "error", err)
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewPlaybarSvcClient(cc)

	for i := range requests {
		_, err := cl.ExecutePlaylistGroupAction(context.Background(), requests[i])
		onerror.Log(err)
	}

	applyplaylistgroupactionschanges()
}

// ExecutePlaybarAction -.
func ExecutePlaybarAction(req *m3uetcpb.ExecutePlaybarActionRequest) (err error) {
	cc, err := getClientConn1()
	if err != nil {
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	_, err = cl.ExecutePlaybarAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		slog.Error(s.Message())
		return
	}
	return
}

// ExecutePlaylistAction -.
func ExecutePlaylistAction(req *m3uetcpb.ExecutePlaylistActionRequest) (
	*m3uetcpb.ExecutePlaylistActionResponse, error) {

	cc, err := getClientConn1()
	if err != nil {
		return nil, err
	}
	defer cc.Close()
	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	res, err := cl.ExecutePlaylistAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return nil, err
	}
	return res, nil
}

// ExecutePlaylistGroupAction -.
func ExecutePlaylistGroupAction(req *m3uetcpb.ExecutePlaylistGroupActionRequest) (
	*m3uetcpb.ExecutePlaylistGroupActionResponse, error) {

	cc, err := getClientConn1()
	if err != nil {
		return nil, err
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	res, err := cl.ExecutePlaylistGroupAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return nil, err
	}
	return res, nil
}

// ExecutePlaylistTrackAction -.
func ExecutePlaylistTrackAction(req *m3uetcpb.ExecutePlaylistTrackActionRequest) (err error) {
	cc, err := getClientConn1()
	if err != nil {
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	_, err = cl.ExecutePlaylistTrackAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		slog.Error(s.Message())
		return
	}
	return
}

// ExportPlaylist -.
func ExportPlaylist(req *m3uetcpb.ExportPlaylistRequest) (err error) {
	cc, err := getClientConn1()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	_, err = cl.ExportPlaylist(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}
	return
}

// ImportPlaylists -.
func ImportPlaylists(req *m3uetcpb.ImportPlaylistsRequest) (
	msgList []string, err error) {

	cc, err := getClientConn1()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	stream, err := cl.ImportPlaylists(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		slog.Error(s.Message())
		return
	}

	for {
		var res *m3uetcpb.ImportPlaylistsResponse
		res, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		msgList = append(msgList, res.ImportErrors...)
	}
	return
}

func applyplaylistgroupactionschanges() {
	slog.Info("Applying playlist group actions changes")

	toRemove := store.BData.PlaylistGroupActionsChanges()

	cc, err := getClientConn1()
	if err != nil {
		slog.Error("Failed to get client connection", "error", err)
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewPlaybarSvcClient(cc)

	for _, id := range toRemove {
		req := &m3uetcpb.ExecutePlaylistGroupActionRequest{
			Action: m3uetcpb.PlaylistGroupAction_PG_DESTROY,
			Id:     id,
		}
		_, err := cl.ExecutePlaylistGroupAction(context.Background(), req)
		onerror.Log(err)
	}
}

func subscribeToPlaybarStore() {
	slog.Info("Subscribing to playbar store")

	defer wgplaybar.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		slog.Error("Failed to obtain client connection", "error", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	stream, err := cl.SubscribeToPlaybarStore(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		slog.Error("Failed to subscribe to playbar store", "error", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			slog.Info("Subscription closed by server", "error", err)
			break
		}

		store.BData.ProcessSubscriptionResponse(res)

		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}

}

func unsubscribeFromPlaybarStore() {
	slog.Info("Unsubscribing from playbar store")

	id := store.BData.SubscriptionID()

	cc, err := getClientConn()
	if err != nil {
		slog.Error("Failed to obtain client connection", "error", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	_, err = cl.UnsubscribeFromPlaybarStore(
		context.Background(),
		&m3uetcpb.UnsubscribeFromPlaybarStoreRequest{
			SubscriptionId: id,
		},
	)
	onerror.Log(err)
}

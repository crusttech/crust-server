// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/messaging/types"
)

func TestChannel(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := Channel(context.Background(), factory.Database.MustGet())
	chn := &types.Channel{}

	var name1, name2 = "Test channel v1", "Test channel v2"

	var cc []*types.Channel

	{
		chn.Name = name1
		chn, err = rpo.CreateChannel(chn)
		test.Assert(t, err == nil, "CreateChannel error: %+v", err)
		test.Assert(t, chn.Name == name1, "Changes were not stored")

		{
			chn.Name = name2

			chn, err = rpo.UpdateChannel(chn)
			test.Assert(t, err == nil, "UpdateChannel error: %+v", err)
			test.Assert(t, chn.Name == name2, "Changes were not stored")
		}

		{
			chn, err = rpo.FindChannelByID(chn.ID)
			test.Assert(t, err == nil, "FindChannelByID error: %+v", err)
			test.Assert(t, chn.Name == name2, "Changes were not stored")
		}

		{
			cc, err = rpo.FindChannels(&types.ChannelFilter{Query: name2})
			test.Assert(t, err == nil, "FindChannels error: %+v", err)
			test.Assert(t, len(cc) > 0, "No results found")
		}

		{
			err = rpo.ArchiveChannelByID(chn.ID)
			test.Assert(t, err == nil, "ArchiveChannelByID error: %+v", err)
		}

		{
			err = rpo.UnarchiveChannelByID(chn.ID)
			test.Assert(t, err == nil, "UnarchiveChannelByID error: %+v", err)
		}

		{
			err = rpo.DeleteChannelByID(chn.ID)
			test.Assert(t, err == nil, "DeleteChannelByID error: %+v", err)
		}

		{
			err = rpo.UndeleteChannelByID(chn.ID)
			test.Assert(t, err == nil, "UndeleteChannelByID error: %+v", err)
		}
	}
}

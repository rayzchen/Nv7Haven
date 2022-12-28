package categories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

// contains is for whether to make sure the cat contains the element (needed for rmcat)
func (c *Categories) catEditCmd(ctx sevcord.Ctx, cat string, elem int, kind types.PollKind, format string, contains bool) {
	ctx.Acknowledge()

	// Get category
	var name string
	err := c.db.QueryRow(`SELECT name FROM categories WHERE guild=$1 AND LOWER(name)=$2`, ctx.Guild(), strings.ToLower(cat)).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows && !contains {
			name = cat
		} else {
			c.base.Error(ctx, err)
			return
		}
	}

	// Get element
	elemName, err := c.base.GetName(ctx.Guild(), elem)
	if err != nil {
		c.base.Error(ctx, err)
		return
	}

	// Check if contains
	if contains {
		var cont bool
		err := c.db.QueryRow(`SELECT $1=ANY(elements) FROM categories WHERE guild=$2 AND name=$3`, elem, ctx.Guild(), name).Scan(&cont)
		if err != nil {
			c.base.Error(ctx, err)
			return
		}
		if !cont {
			ctx.Respond(sevcord.NewMessage(fmt.Sprintf("Element **%s** is not in category **%s**! "+types.RedCircle, elemName, name)))
			return
		}
	}

	// Make poll
	err = c.polls.CreatePoll(ctx, &types.Poll{
		Kind: kind,
		Data: types.PgData{
			"cat":   name,
			"elems": []any{any(float64(elem))},
		},
	})
	if err != nil {
		c.base.Error(ctx, err)
		return
	}

	// Respond
	ctx.Respond(sevcord.NewMessage(fmt.Sprintf(format, elemName, name)))
}

func (c *Categories) AddCat(ctx sevcord.Ctx, opts []any) {
	c.catEditCmd(ctx, opts[0].(string), int(opts[1].(int64)), types.PollKindCategorize, "Suggested to add **%s** to **%s** 🗃️", false)
}

func (c *Categories) RmCat(ctx sevcord.Ctx, opts []any) {
	c.catEditCmd(ctx, opts[0].(string), int(opts[1].(int64)), types.PollKindUncategorize, "Suggested to remove **%s** from **%s** 🗃️", true)
}

func (c *Categories) DelCat(ctx sevcord.Ctx, opts []any) {
	ctx.Acknowledge()

	// Get category
	var name string
	var els pq.Int32Array
	err := c.db.QueryRow(`SELECT name, elements FROM categories WHERE guild=$1 AND LOWER(name)=$2`, ctx.Guild(), strings.ToLower(opts[0].(string))).Scan(&name, &els)
	if err != nil {
		c.base.Error(ctx, err)
		return
	}

	// Make poll
	err = c.polls.CreatePoll(ctx, &types.Poll{
		Kind: types.PollKindUncategorize,
		Data: types.PgData{
			"cat":   name,
			"elems": util.Map([]int32(els), func(v int32) any { return any(float64(v)) }),
		},
	})
	if err != nil {
		c.base.Error(ctx, err)
		return
	}

	// Respond
	ctx.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested to delete category **%s** 🗃️", name)))
}

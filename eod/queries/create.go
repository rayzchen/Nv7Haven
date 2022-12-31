package queries

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
)

func (q *Queries) createCmd(c sevcord.Ctx, name string, kind types.QueryKind, data map[string]any) {
	c.Acknowledge()

	// Check if name exists
	var edit bool
	err := q.db.Get(&edit, "SELECT EXISTS(SELECT 1 FROM queries WHERE LOWER(name)=$1 AND guild=$2)", strings.ToLower(name), c.Guild())
	if err != nil {
		q.base.Error(c, err)
		return
	}
	if edit {
		err = q.db.QueryRow(`SELECT name FROM queries WHERE LOWER(name)=$1 AND guild=$2`, strings.ToLower(name), c.Guild()).Scan(&name)
		if err != nil {
			q.base.Error(c, err)
			return
		}
	} else {
		// Fix name
		var ok types.Resp
		name, ok = base.CheckName(name)
		if !ok.Ok {
			c.Respond(sevcord.NewMessage(ok.Message + " " + types.RedCircle))
			return
		}
	}

	// Check if data already exists
	var existsName string
	err = q.db.QueryRow("SELECT name FROM queries WHERE data@>$1 AND data<@$1 AND kind=$3 AND guild=$2", types.PgData(data), c.Guild(), string(kind)).Scan(&existsName)
	if err != nil && err != sql.ErrNoRows {
		q.base.Error(c, err)
		return
	}
	if err == nil {
		c.Respond(sevcord.NewMessage(fmt.Sprintf("Query **%s** already exists with this data! "+types.RedCircle, existsName)))
		return
	}

	// Create
	err = q.polls.CreatePoll(c, &types.Poll{
		Kind: types.PollKindQuery,
		Data: types.PgData{
			"query": name,
			"edit":  edit,
			"kind":  string(kind),
			"data":  any(data),
		},
	})
	if err != nil {
		q.base.Error(c, err)
		return
	}

	// Respond
	word := "create"
	if edit {
		word = "edit"
	}
	c.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested to %s query! 🧮", word)))
}

func (q *Queries) CreateElementsCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	q.createCmd(c, opts[0].(string), types.QueryKindElements, map[string]any{})
}

func (q *Queries) CreateElementCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	// Check if element exists
	var exists bool
	err := q.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM elements WHERE id=$1)", opts[1].(int64))
	if err != nil {
		q.base.Error(c, err)
		return
	}
	if !exists {
		c.Respond(sevcord.NewMessage("Element does not exist! " + types.RedCircle))
		return
	}
	q.createCmd(c, opts[0].(string), types.QueryKindElement, map[string]any{"elem": float64(opts[1].(int64))})
}

func (q *Queries) CreateCategoryCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	// Get name
	var name string
	err := q.db.Get(&name, "SELECT name FROM categories WHERE LOWER(name)=$1", strings.ToLower(opts[1].(string)))
	if err != nil {
		q.base.Error(c, err)
		return
	}
	q.createCmd(c, opts[0].(string), types.QueryKindCategory, map[string]any{"cat": name})
}

func (q *Queries) CreateProductsCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	// Get name
	var name string
	err := q.db.Get(&name, "SELECT name FROM queries WHERE LOWER(name)=$1", strings.ToLower(opts[1].(string)))
	if err != nil {
		q.base.Error(c, err)
		return
	}
	q.createCmd(c, opts[0].(string), types.QueryKindProducts, map[string]any{"query": name})
}

func (q *Queries) CreateParentsCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	// Get name
	var name string
	err := q.db.Get(&name, "SELECT name FROM queries WHERE LOWER(name)=$1", strings.ToLower(opts[1].(string)))
	if err != nil {
		q.base.Error(c, err)
		return
	}
	q.createCmd(c, opts[0].(string), types.QueryKindParents, map[string]any{"query": name})
}

func (q *Queries) CreateInventoryCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	q.createCmd(c, opts[0].(string), types.QueryKindInventory, map[string]any{"user": opts[1].(*discordgo.User).ID})
}

func (q *Queries) CreateRegexCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	// Get name
	var name string
	err := q.db.Get(&name, "SELECT name FROM queries WHERE LOWER(name)=$1", strings.ToLower(opts[1].(string)))
	if err != nil {
		q.base.Error(c, err)
		return
	}
	// Check regex
	_, err = regexp.CompilePOSIX(opts[2].(string))
	if err != nil {
		q.base.Error(c, err)
		return
	}
	q.createCmd(c, opts[0].(string), types.QueryKindRegex, map[string]any{"query": name, "regex": opts[2].(string)})
}

func (q *Queries) CreateComparisonCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	// Get name
	var name string
	err := q.db.Get(&name, "SELECT name FROM queries WHERE LOWER(name)=$1", strings.ToLower(opts[1].(string)))
	if err != nil {
		q.base.Error(c, err)
		return
	}
	// Parse if needed
	val := any(opts[4].(string))
	switch opts[2].(string) {
	case "treesize", "color":
		intV, err := strconv.Atoi(opts[4].(string))
		if err != nil {
			q.base.Error(c, err)
			return
		}
		val = any(float64(intV))
	}
	q.createCmd(c, opts[0].(string), types.QueryKindComparison, map[string]any{"query": name, "field": opts[2].(string), "typ": opts[3].(string), "value": val})
}

func (q *Queries) CreateOperationCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	// Get names
	var nameLeft string
	err := q.db.Get(&nameLeft, "SELECT name FROM queries WHERE LOWER(name)=$1", strings.ToLower(opts[1].(string)))
	if err != nil {
		q.base.Error(c, err)
		return
	}
	var nameRight string
	err = q.db.Get(&nameRight, "SELECT name FROM queries WHERE LOWER(name)=$1", strings.ToLower(opts[2].(string)))
	if err != nil {
		q.base.Error(c, err)
		return
	}
	q.createCmd(c, opts[0].(string), types.QueryKindOperation, map[string]any{"left": nameLeft, "right": nameRight, "op": opts[3].(string)})
}

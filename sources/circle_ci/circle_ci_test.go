package circle_ci

import (
	"testing"

	"github.com/andreiko/alfred-sources/sources"
	"github.com/stretchr/testify/assert"
)

func TestNewItem(t *testing.T) {
	item := NewCircleCiItem("github", "dotfiles", "andreiko", "https://circleci.com")
	assert.Equal(t, "andreiko/dotfiles", item.Autocomplete())
}

func TestAttributes(t *testing.T) {
	expected := map[string]interface{}{
		"vs_type":  "github",
		"reponame": "dotfiles",
		"username": "andreiko",
		"fullname": "andreiko/dotfiles",
	}

	item := NewCircleCiItem("github", "dotfiles", "andreiko", "https://circleci.com")
	assert.Equal(t, expected, item.Attributes())
}

func TestRanks(t *testing.T) {
	item := NewCircleCiItem("github", "dotfiles", "andreiko", "https://circleci.com")

	fullname_match := item.GetRank("andreiko/dotfiles")
	reponame_match := item.GetRank("dotfiles")
	fullname_beginning_match := item.GetRank("andreiko/do")
	reponame_beginning_match := item.GetRank("dot")
	reponame_substring_match := item.GetRank("files")
	username_substring_match := item.GetRank("eik")

	assert.True(t, fullname_match > reponame_match)
	assert.True(t, reponame_match > fullname_beginning_match)
	assert.True(t, fullname_beginning_match > reponame_beginning_match)
	assert.True(t, reponame_beginning_match > reponame_substring_match)
	assert.True(t, reponame_substring_match > username_substring_match)

	no_match := item.GetRank("x")
	assert.Equal(t, 0, no_match)
}

func TestLessThan(t *testing.T) {
	i1 := NewCircleCiItem("github", "abc1", "andreiko", "https://circleci.com")
	i2 := NewCircleCiItem("github", "abc2", "andreiko", "https://circleci.com")
	i3 := NewCircleCiItem("github", "abc3", "andreiko", "https://circleci.com")

	assert.True(t, i1.LessThan(i2))
	assert.True(t, i2.LessThan(i3))
}

func testItems() []sources.Item {
	return []sources.Item{
		NewCircleCiItem("github", "abc3", "andreiko", "https://circleci.com"),
		NewCircleCiItem("github", "abc2", "andreiko", "https://circleci.com"),
		NewCircleCiItem("github", "bc", "andreiko", "https://circleci.com"),
		NewCircleCiItem("github", "xbnc", "andreiko", "https://circleci.com"),
	}
}

func TestQuery(t *testing.T) {
	src := new(CircleCiSource)
	src.items = testItems()

	result := src.Query("bc")
	assert.Equal(t, 3, len(result))

	assert.Equal(t, "andreiko/bc", result[0].Autocomplete())
	assert.Equal(t, "andreiko/abc2", result[1].Autocomplete())
	assert.Equal(t, "andreiko/abc3", result[2].Autocomplete())
}

func TestLongQuery(t *testing.T) {
	src := new(CircleCiSource)
	src.items = testItems()

	result := src.Query("abcdefghijklmnopqrstuvwxyz")
	assert.Equal(t, 0, len(result))
}

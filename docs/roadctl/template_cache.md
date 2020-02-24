# Template cache

Starting with version 0.6.0 templates are checked out using
git clone and the life cycle events affecting that cache
are stored in .pr_cache in the template directory.

Code pertaining to caches is stored in cmd/template_cache.go.

In prior versions, the templates locations

```Go
var defaultTemplateDir = "templates"
```

As of 0.6.0, you must create a tplCache object using NewTemplateCache.
This validates the location and sets **defaultTemplateDir** correctly.

The defaultTemplateDir is being depreciated.

```go
	tc, err := NewTemplateCache()

	// Use this
	tc.location.Location()

	// not this
	defaultTemplateDir
```

## New defaults

$HOME/.pavedroad.d/templates is the new default

You can change the default with:

- export PR_TEMPLATE_DIR="some new location"

- roadctl cmd resource --templates "some new location"

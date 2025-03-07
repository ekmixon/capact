package validation

import (
	"context"
	"fmt"
	"strings"

	"capact.io/capact/internal/multierror"
	"capact.io/capact/internal/ptr"
	gqlpublicapi "capact.io/capact/pkg/hub/api/graphql/public"
	"github.com/pkg/errors"
)

// HubClient defines Hub methods needed for ResolveTypeRefsToJSONSchemas.
type HubClient interface {
	ListTypeRefRevisionsJSONSchemas(ctx context.Context, filter gqlpublicapi.TypeFilter) ([]*gqlpublicapi.TypeRevision, error)
}

// ResolveTypeRefsToJSONSchemas resolves Type references to theirs JSON schemas.
func ResolveTypeRefsToJSONSchemas(ctx context.Context, hubCli HubClient, inTypeRefs TypeRefCollection) (SchemaCollection, error) {
	// 1. Fetch revisions for given TypeRefs
	var typeRefsPath []string
	for _, ref := range inTypeRefs {
		typeRefsPath = append(typeRefsPath, ref.Path)
	}
	// No TypeRefs that should be resolved, early return to do not call Hub
	if len(typeRefsPath) == 0 {
		return nil, nil
	}

	typeRefsPathFilter := fmt.Sprintf(`(%s)`, strings.Join(typeRefsPath, "|"))
	gotTypes, err := hubCli.ListTypeRefRevisionsJSONSchemas(ctx, gqlpublicapi.TypeFilter{
		PathPattern: ptr.String(typeRefsPathFilter),
	})
	if err != nil {
		return nil, errors.Wrap(err, "while fetching JSONSchemas for input TypeRefs")
	}

	indexedTypes := map[string]interface{}{}
	for _, rev := range gotTypes {
		if rev == nil || rev.Spec == nil {
			continue
		}
		key := fmt.Sprintf("%s:%s", rev.Metadata.Path, rev.Revision)
		indexedTypes[key] = rev.Spec.JSONSchema
	}

	var (
		merr    = multierror.New()
		schemas = SchemaCollection{}
	)
	for name, ref := range inTypeRefs {
		refKey := fmt.Sprintf("%s:%s", ref.Path, ref.Revision)
		schema, found := indexedTypes[refKey]
		if !found {
			// It means that manifest refers to not existing TypeRef.
			// From our perspective it's error - we shouldn't have invalid manifests in Hub.
			merr = multierror.Append(merr, fmt.Errorf("TypeRef %q was not found in Hub", refKey))
			continue
		}
		str, ok := schema.(string)
		if !ok {
			merr = multierror.Append(merr, fmt.Errorf("unexpected JSONSchema type for %q: expected %T, got %T", refKey, "", schema))
			continue
		}
		schemas[name] = Schema{
			Value:    str,
			Required: ref.Required,
		}
	}

	if err := merr.ErrorOrNil(); err != nil {
		return nil, err
	}

	return schemas, nil
}

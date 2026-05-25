package app

import (
	"io"
	"log/slog"

	"github.com/marvinlanhenke/terraview/internal/planview"
)

// discardLogger returns a no-op logger suitable for tests.
func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// testPlanRoot returns a minimal plan node tree for use in tests.
func testPlanRoot() *planview.Node {
	return &planview.Node{
		Id:     "root",
		Label:  "Root",
		Kind:   planview.NodeGroup,
		Action: planview.ActionNoOp,
		Children: []*planview.Node{
			{
				Id:         "create-group",
				Label:      "Create",
				LabelCount: "(1/2)",
				Kind:       planview.NodeGroup,
				Action:     planview.ActionCreate,
				Children: []*planview.Node{
					{
						Id:      "aws_instance.web",
						Label:   "aws_instance.web",
						Kind:    planview.NodeResource,
						Action:  planview.ActionCreate,
						Payload: map[string]any{"address": "aws_instance.web"},
					},
				},
			},
			{
				Id:         "delete-group",
				Label:      "Delete",
				LabelCount: "(1/2)",
				Kind:       planview.NodeGroup,
				Action:     planview.ActionDelete,
				Children: []*planview.Node{
					{
						Id:      "aws_s3_bucket.old",
						Label:   "aws_s3_bucket.old",
						Kind:    planview.NodeResource,
						Action:  planview.ActionDelete,
						Payload: map[string]any{"address": "aws_s3_bucket.old"},
					},
				},
			},
		},
	}
}

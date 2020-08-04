package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/slack-go/slack"
	cbpb "google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
)

func TestWriteMessage(t *testing.T) {
	n := new(slackNotifier)

	b := &cbpb.Build{
		ProjectId: "my-project-id",
		Status:    cbpb.Build_SUCCESS,
		LogUrl:    "https://some.example.com/log/url?foo=bar",
	}

	b.Substitutions = make(map[string]string)
	b.Substitutions["REPO_NAME"] = "some-repo-name"
	b.Substitutions["BRANCH_NAME"] = "branch"
	b.Substitutions["SHORT_SHA"] = "1234"


	got, err := n.writeMessage(b)
	if err != nil {
		t.Fatalf("writeMessage failed: %v", err)
	}

	want := &slack.WebhookMessage{
		Attachments: []slack.Attachment{{
			Text:  "some-repo-name, branch, 1234: SUCCESS",
			Color: "good",
			Actions: []slack.AttachmentAction{{
				Text: "View Logs",
				Type: "button",
				URL:  "https://some.example.com/log/url?foo=bar&utm_campaign=google-cloud-build-notifiers&utm_medium=chat&utm_source=google-cloud-build",
			}},
		}},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("writeMessage got unexpected diff: %s", diff)
	}
}

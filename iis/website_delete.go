package iis

import "context"

func (client Client) DeleteWebsite(ctx context.Context, id string) error {
	return httpDelete(ctx, client, client.websitePath(id))
}

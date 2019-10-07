package opsmanager

func (client opsManagerClient) DeleteProject(projectID string) error {
	url := client.resolver.Of("/groups/%s", projectID)
	return client.Delete(url)
}

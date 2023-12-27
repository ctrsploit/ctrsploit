package version

type Map map[string]Version

func (m Map) Get(keys []string) (versions []Version) {
	for _, key := range keys {
		versions = append(versions, m[key])
	}
	return
}

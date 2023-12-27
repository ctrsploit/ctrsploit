package version

// TODO: ordered map

type Map map[string]Version

func (m Map) Get(keys []string) (versions []Version) {
	for _, key := range keys {
		versions = append(versions, m[key])
	}
	return
}

func (m Map) Values() (versions []Version) {
	for _, v := range m {
		versions = append(versions, v)
	}
	return
}

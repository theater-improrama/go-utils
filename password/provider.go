package password

type ProviderFn func() (Provider, error)

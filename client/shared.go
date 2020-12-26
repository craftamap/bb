package client

type Link struct {
	Href string `mapstructure:"href"`
}

type Branch struct {
	Name string `mapstructure:"name"`
}

type Account struct {
	AccountID   string `mapstructure:"account_id"`
	DisplayName string `mapstructure:"display_name"`
	Nickname    string `mapstructure:"nickname"`
	Type        string `mapstructure:"user"`
	UUID        string `mapstructure:"uuid"`
}

module github.com/alackey/Stonks

go 1.15

require (
	github.com/alackey/go-tdameritrade v0.0.0-20210113060627-e8520d925a58
	github.com/aws/aws-sdk-go v1.36.19
	github.com/bwmarrin/discordgo v0.22.0
	github.com/dustin/go-humanize v1.0.0
	github.com/joho/godotenv v1.3.0
	github.com/spacecodewor/fmpcloud-go v0.0.0-20201221162404-ee5f5303f8a9
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6
)

replace github.com/zricethezav/go-tdameritrade/tdameritrade => github.com/alackey/go-tdameritrade v0.0.0-20210113060627-e8520d925a58

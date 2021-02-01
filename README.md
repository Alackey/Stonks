<!-- PROJECT LOGO -->
<br />
<p align="center">
  <!--
  <a href="">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>
  -->
  <h1 align="center">Stonks</h1>

  <p align="center">
    A discord bot designed to give users information about the stock market.
    <br />
    <a href="https://github.com/Alackey/Stonks/issues">Report Bug</a>
    ·
    <a href="https://github.com/Alackey/Stonks/issues">Request Feature</a>
  </p>
</p>



<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#acknowledgements">Acknowledgements</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

DISCAIMER: This project is the early stages and the commands can change at any point.

When I looked around, there weren't many Discord bots for the stock market, and the bots I did find weren't in a format that I liked or were missing features I wanted. Since they were lacking things I wanted, I decided to make my own with a nicer format and different features.


### Built With

* [Go](https://golang.org/)
* [DiscordGo](https://github.com/bwmarrin/discordgo)
* [fmpcloud-go](github.com/spacecodewor/fmpcloud-go)


<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple example steps.

### Prerequisites

A couple of things youll need before you can have Narwhal up and running:
* Go 1.x
* [A Discord Bot Token](https://discord.com/developers/applications)
* [A Financial Modeling Prep API key](https://financialmodelingprep.com/developer)

### Installation

1. Get your API keys
2. Clone the repo
   ```sh
   git clone https://github.com/Alackey/Stonks.git
   ```
3. Download the go modules
   ```sh
   go mod download
   ```
4. Set environment variables manually, or copy the .env.example and create a .env file with the environment variables
5. Start the bot
   ```sh
   go run .
   ```


### AWS Elastic Beanstalk
I am currently hosting this bot on AWS Elastic Beanstalk, and to deploy this app to a currently existing Elastic Beakstalk environment you would run the command below. 

```sh
eb deploy {ENVIRONMENT_NAME}
```

This command uses the Buildfile and Procfile for building and deploying the bot to Elastic Beanstalk.


<!-- USAGE EXAMPLES -->
## Usage
All commands begin with the default "$" prefix. For example: 

 ```sh
  $somecommand arg1 arg2 ...arg
  ```

### Commands

**$q \<symbol>** - Gets the price information of a stock based off the symbol

**$futures** - Gets the price information for some futures

**$market** - Shows a heatmap of the market and its sectors

**$market crypto** - Shows a heatmap of the crypto market

**$news \<symbol>** - Gets the most recent news about a stock based off the symbol

**$help** - Shows a list of available commands


<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**. 

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.


<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements
* [Best ReadME Template](https://github.com/othneildrew/Best-README-Template/blob/master/README.md)

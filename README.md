[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/paypi/paypi-go">
    <img src="images/logo.png" alt="Logo" height="80">
  </a>

  <h3 align="center">PayPI Go Client</h3>

  <p align="center">
    Sell your API, today.
    <br />
    <a href="https://partner.paypi.dev/"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://paypi.dev/">Homepage</a>
    ·
    <a href="https://github.com/paypi/paypi-go/issues">Report Bug</a>
    ·
    <a href="https://github.com/paypi/paypi-go/issues">Request Feature</a>
  </p>
</p>

<!-- TABLE OF CONTENTS -->

## Table of Contents

- [About the Project](#about-the-project)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)
- [Acknowledgements](#acknowledgements)

<!-- ABOUT THE PROJECT -->

## About The Project

[![PayPI Screenshot][product-screenshot]](https://paypi.dev)

PayPI makes API creators' lives easier by handling API keys, user accounts, payments and more.
API users have one account to access all APIs using PayPI.

We worry about API authentication and payments so you can focus on making awesome APIs! This library enables you to interact with PayPI from a Golang project.

<!-- GETTING STARTED -->

## Getting Started

> <a href="https://partner.paypi.dev/"><strong>See full documentation here</strong></a>

Use Go Modules to import PayPI:

```go
import (
    "github.com/paypi/paypi-go/paypi"
)
```

Or explicitly fetch with `go get`:

```go
go get github.com/paypi/paypi-go
```

### Initialisation

To setup the library import it and set the API key to the key given to you from the PayPI Dashboard.

```go
import (
  "http"
  "github.com/paypi/paypi-go/paypi"
)

paypi.Key = "<API_SECRET_KEY>"
```

### Authenticate

`Authenticate` is used to check if the user's API key is valid and not blocked via Paypi for any reason (non payment, rate limiting, etc).

```go
user, err := paypi.Authenticate("<USER_TOKEN>")
if err != nil {
  http.Error(w, "User token is unauthorized", http.StatusUnauthorized)
  return
}
```

If this method returns without error you can continue processing the request. If it fails you should immediately return an unauthorized response from your API.

### Make Charges

`user.MakeCharge` is used to charge users on the platform. When `user.MakeCharge` is called the user
is charged for their usage on your platform. This should be done at the end of your request, when you
are sure that their request is going to be fulfilled successfully.

#### Static Charges

Static charges need only provide the charge identifier

```go
user, _ := paypi.Authenticate("<USER_TOKEN>")

charge, err := user.MakeCharge(paypi.MakeChargeInput{
  ChargeIdentifier: "<STATIC_CHARGE_IDENTIFIER>",
})
```

#### Dynamic Charges

If a charge is dynamic, you must provide the number of units used. This is usually representative of
seconds of processing time, or MB of data processed. This unit is multiplied by the price set for the
dynamic charge to calculate the total cost to the user.

```go
user, _ := paypi.Authenticate("<USER_TOKEN>")

charge, err := user.MakeCharge(paypi.MakeChargeInput{
  ChargeIdentifier: "<STATIC_CHARGE_IDENTIFIER>",
  UnitsUsed: 5,
})
```

## Example Usage

```go
package main

import (
  "http"
  "github.com/paypi/paypi-go/paypi"
)

// Set the API secret key on initialising your application
paypi.Key = "<API_SECRET_KEY>"

func handleRequest(w http.ResponseWriter, r *http.Request) {
  // Check if a user's token is valid, this token usually comes from the
  // Authorization header, but can be given in any way you choose.
  user, err := paypi.Authenticate("<USER_TOKEN>")
  if err != nil {
    http.Error(w, "User token is unauthorized", http.StatusUnauthorized)
    return
  }

  // Do some processing, fetch the response etc...

  // If request was successful, make a static charge
  charge, err := user.MakeCharge(paypi.MakeChargeInput{
    ChargeIdentifier: "<STATIC_CHARGE_IDENTIFIER>",
  })

  // Make a dynamic charge based on processing
  charge, err = user.MakeCharge(paypi.MakeChargeInput{
    ChargeIdentifier: "<DYNAMIC_CHARGE_IDENTIFIER>",
    UnitsUsed: 2.34,
  })

  fmt.Fprintf(w, "Request Successful")
}
```

<!-- ROADMAP -->

## Roadmap

See the [open issues](https://github.com/paypi/paypi-go/issues) for a list of proposed features (and known issues).

<!-- CONTRIBUTING -->

## Contributing

All contributions are welcome. Please follow this workflow:

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- LICENSE -->

## License

All rights reserved.

<!-- CONTACT -->

## Contact

Alex - alex@paypi.dev  
Tom - tom@paypi.dev

Project Link: [https://github.com/paypi/paypi-go](https://github.com/paypi/paypi-go)

<!-- ACKNOWLEDGEMENTS -->

## Acknowledgements

- [Img Shields](https://shields.io)
- [Choose an Open Source License](https://choosealicense.com)

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[contributors-shield]: https://img.shields.io/github/contributors/Paypi/paypi-go.svg?style=flat-square
[contributors-url]: https://github.com/paypi/paypi-go/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/Paypi/paypi-go.svg?style=flat-square
[forks-url]: https://github.com/paypi/paypi-go/network/members
[stars-shield]: https://img.shields.io/github/stars/Paypi/paypi-go.svg?style=flat-square
[stars-url]: https://github.com/paypi/paypi-go/stargazers
[issues-shield]: https://img.shields.io/github/issues/Paypi/paypi-go.svg?style=flat-square
[issues-url]: https://github.com/paypi/paypi-go/issues
[license-shield]: https://img.shields.io/github/license/Paypi/paypi-go.svg?style=flat-square
[license-url]: https://github.com/paypi/paypi-go/blob/master/LICENSE.txt
[product-screenshot]: images/product.png

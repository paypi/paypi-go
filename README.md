# paypi-go

The official Paypi Go Merchant client library

## Setup
Use gomodules to import:
```go
import (
    "github.com/Paypi/paypi-go/paypi"
)
```

or explicitly fetch with `go get`

```go
go get github.com/Paypi/paypi-go
```


## Documentation
The Paypi library has two main methods, `Authenticate` and `MakeCharge`.

### Initialisation
To setup the library just import it and set the API secret key for the API you are implementing.
This needs to be done before you call any other methods. You can find this key on the Paypi website
under the API you are implementing Paypi with.

```go 
import (
  "http"
  "github.com/Paypi/paypi-go/paypi"
)

paypi.Key = "<API_SECRET_KEY>"
```

### Authenticate
`Authenticate` is used to check if the user's API key is valid and not blocked via Paypi 
for any reason (non payment, rate limiting, etc). 

```go 
user, err := paypi.Authenticate("<USER_TOKEN>")
if err != nil {
  http.Error(w, "User token is unauthorized", http.StatusUnauthorized)
  return
}
```

If this method returns without error you can continue processing the request. If it fails
you should immediately return an unauthorized response from your API.

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
  "github.com/Paypi/paypi-go/paypi"
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

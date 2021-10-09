

<p align="center">
  <!-- <a href="https://github.com/devashar13">
    
  </a> -->

  <h3 align="center">Instagram API</h3>



<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

This is an example of how to list things you need to use the software and how to install them.
* go


### Installation
 
1. Clone the Repo
```sh
go get github.com/devashar13/instagram-api
```

## Routes


```sh
/users
```
* Query Parameters
```sh
POST
{
    "name":"Name",
    "email":"name@gmail.com",
    "password":"password"
}
```
Password Hashing Algorithm Used: SHA256

```sh
GET
/users/<userid>
```

```sh
POST
/products
```

```sh
POST
{
    "user":"userid",
    "caption":"caption",
    "imageurl":"imageurl"
    
}
```
Time is set to the time of upload 

```sh
GET
/products/<productid>
```

```sh
GET
/products/users/<userid>
```
or To paginate
```sh
GET
/products/users/<userid>?limit=2
```
##  Testing

```
go test
```

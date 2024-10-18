# Payment and Login Transactions in Go

This repository contains various payment and login integrations implemented in Go, including support for Google Login, Google Maps integration, PayPal payments, Adyen payments, invoice generation, and password management.

## Features

### 1. Google Login
- Implemented using OAuth2 for authentication and authorization
- Allows users to log in with their Google accounts

### 2. Google Maps Integration
- Provides integration with the Google Maps API for services such as distance calculation, geocoding, and displaying maps
- Can be used to enhance location-based features in applications

### 3. PayPal Payments
- Supports payment processing using the PayPal REST API
- Allows for payment authorization, capture, and refund functionalities

### 4. Adyen Payments
- Integration with Adyen for secure payment processing
- Supports multiple payment methods and currencies
- Includes features for payment authorization and capturing

### 5. Invoice Generation
- Supports creating invoices in PDF format
- Can be customized to include details such as items, prices, taxes, and company information
- Suitable for generating standard-compliant invoices for various regions

### 6. Password Management
- Includes functions for securely creating and storing hashed passwords
- Supports password validation and encryption techniques for enhanced security

## Getting Started

Follow these steps to get started with this repository:

### Prerequisites
- Go (version 1.18 or later)
- A valid Google API key for Google Login and Google Maps integrations
- PayPal and Adyen sandbox/test account credentials for payment integration

### Installation

1. Clone the repository:
```bash
git clone https://github.com/SadikSunbul/Payment-and-login-transactions-written-in-go.git
```

2. Navigate to the project directory:
```bash
cd Payment-and-login-transactions-written-in-go
```

3. Install the required Go dependencies:
```bash
go mod tidy
```

### Configuration

1. **Google Login and Google Maps**
   - Set up your Google API credentials and update the configuration file accordingly

2. **PayPal Integration**
   - Add your PayPal REST API credentials to the configuration file

3. **Adyen Integration**
   - Include your Adyen API credentials and configure the payment settings

4. **Invoice Generation**
   - Configure invoice settings such as company details, tax rates, and output folder

### Running the Application

To run the application, use the following command:
```bash
go run main.go
```

The server should start, and you can access the services via the defined endpoints.

## To-Do List
- Apple Login
- Apple Pay
- Google Pay
- Unit tests (optional)

## Contributing

Feel free to submit issues, fork the repository, and send pull requests if you'd like to contribute to this project.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

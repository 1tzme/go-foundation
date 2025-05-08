# 💳 creditcard CLI

## 📦 Features

- `validate` — Validates a credit card number using the Luhn algorithm.
- `generate` — Generates valid card numbers based on a given template.
- `information` — Detects the brand and issuer of a card number using external files.
- `issue` — Issues a new card number based on a specified brand and issuer.

## Installation

```bash
git clone https://github.com/zaaripzha/creditcard.git
cd creditcard
go build -o creditcard .
```

## Usage
Validate a card number
```bash
./creditcard validate 1234567890123456
```
Or from standard input:
```bash
echo "4400430180300003" | ./creditcard validate --stdin
```
Generate card numbers
Using a template (up to 4 trailing asterisks *):
```bash
./creditcard generate "440043018030****"
```
Generate a single random valid card number:
```bash
./creditcard generate --pick "440043018030****"
```
Get card information
```bash
./creditcard information --brands=brands.txt --issuers=issuers.txt "4400430180300003"
```
Issue a card by brand and issuer
```bash
./creditcard issue --brands=brands.txt --issuers=issuers.txt --brand=VISA --issuer=CHAS
```
## Format of brands.txt and issuers.txt
### Each line should follow the format:

```txt
NAME:PREFIX
```
### Example:
```txt
VISA:4
MASTERCARD:51
MASTERCARD:52
MASTERCARD:53
MASTERCARD:54
MASTERCARD:55
AMEX:34
AMEX:37
```
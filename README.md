# Form3 Take Home Exercise
This client present an interface to communicate with the fake api account client from form3. It consists of three methods`fetch an account`, `create an account` and `delete an account`.

## How to run it
This client can be run using `docker-compose.yml` present in the project folder.

on terminal hit:
- docker-compose up

### let's create an account first
To create an account, open postman and select `post` method, in url use -> `http://localhost:10000/form3Client/accounts`.
The post body should look like :
```json
{
  "data": {
    "type": "accounts",
    "id": "cb1e2074-1056-4b27-b4e0-ed9f0c46b066",
    "organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
    "attributes": {
      "country": "GB",
      "base_currency": "GBP",
      "bank_id": "400300",
      "bank_id_code": "GBDSC",
      "bic": "NWBKGB22",
      "name": [
        "Samantha Holder"
      ],
      "alternative_names": [
        "Sam Holder"
      ],
      "account_classification": "Personal",
      "joint_account": false,
      "account_matching_opt_out": false,
      "secondary_identification": "A1B2C3D4"
    }
  }
}
```
Once you created an account, it is time to retrieve it. To retrieve an account, in postman use `GET` method and hit the 
`http://localhost:10000/form3Client/accounts/cb1e2074-1056-4b27-b4e0-ed9f0c46b066` end-point
You should receive the account you created in previous step.

Now it is time to delete the account created, whenever you create an account it gets a **version** as well. To delete an account, we need to have accountId and version otherwise, we won't be able to delete it.
Do delete select `DELETE` as a method in postman and hit -> `http://localhost:10000/form3Client/accounts/cb1e2074-1056-4b27-b4e0-ed9f0c46b066?version=0`

***Voila***, we tested all the happy path of our client

## Improvements

- Better error handling, a robost error struct with more validation on input data would have reduced the
amount of calls we made to fake api account.
- We could have used `Go playground validator` or `custom validator` to verify our input json.

## Developed By
with love
**[Shams Abubakar Azad](https://www.linkedin.com/in/shamsazad/)**
---
title: "Task 3 - Schema Protection"
linkTitle: "Task 3 - Schema Protection"
weight: 30
---

#### Open API Validation/Schema protection

In this task, lets run through the open API/Swagger based schema protection with Fortiweb cloud. Swagger, now known as the OpenAPI Specification (OAS), is a framework for API development that allows developers to design, build, document, and consume RESTful web services.

example of Swagger: https://petstore.swagger.io/

FortiWeb can validate incoming requests against your OpenAPI schema to ensure they conform to the defined structure and data types. This helps prevent injection attacks and other malicious activities.

1. Download the juiceshop schema file to your local machine by clicking on URL below.


https://juiceshopswagger.blob.core.windows.net/juiceshopswagger/swagger.yaml?sp=r&st=2024-08-06T16:05:20Z&se=2024-11-09T01:05:20Z&spr=https&sv=2022-11-02&sr=b&sig=F8TWuKSH430782%2FgJBWLhCQuEDK2101CChRkXx4XdU0%3D


2. From the FortiWeb Cloud Console left pane, select ADD MODULES. Scroll down and turn on  under API Protection to add OPEN API VALIDATION

![apischema1](api-schema1.png)

3. In the API protection module, click on Open API validation > Create OpenAPI Validation Rule. 

![apischema2](api-schema2.png)

4. Click on "choose file" to upload the file downloaded in Step 1, Click OK. 

![apischema3](api-schema3.png)

5. Dont forget to Save at the bottom. 

![apischema4](api-schema4.png)

{{% notice warning %}}
If for some reason you are logged out when you click save here, you will need to log back in using this link ```https://customersso1.fortinet.com/saml-idp/proxy/demo_sallam_okta/login``` and the credentials recieved in the original email.  You will need to repeat steps 1 through 5.
{{%/notice%}}

6. Now lets open POSTMAN, by opening a new terminal (not bash) and type ```Postman``` at the prompt.  This should start the postman application.

- When postman comes up, select "Or continue with lightweight API client"

![postmanlite](p-light.png)

7. we will send a POST request to the URL we have documented in Schema. Change "**GET**" to "**POST**", for URL use: **https://yournumber.fwebtraincse.com/b2b/v2/orders**

for Request body, Click on Body > Raw > JSON and paste the following:

```sh
{
  "cid": "testing",
  "orderLines": [
    {
      "productId": "testing",
      "quantity": 500,
      "customerReference": 1
    }
  ],
  "orderLinesData": "[{\"productId\": 12,\"quantity\": 10000,\"customerReference\": [\"PO0000001.2\", \"SM20180105|042\"],\"couponCode\": \"pes[Bh.u*t\"},{\"productId\": 13,\"quantity\": 2000,\"customerReference\": \"PO0000003.4\"}]"
}
```
![apischema6](api-schema6.png)

- Note: The schema for Product ID is changed from Integer to String. the Fortiweb cloud Juiceshop schema we uploaded have this value defined as Integer. 

![apischema10](api-schema10.png)

- Click on "**SEND**"

8. We will see "403 internal server error" with a Fortiweb cloud block message in HTML.

![apischema7](api-schema7.png)

9. In Fortiweb Cloud on the left hand side of the screen go to Threat Analytics > Attack log > we can see a log generated for this block request to show the reason for block is Open API schema Violation. 

![apischema8](api-schema8.png)

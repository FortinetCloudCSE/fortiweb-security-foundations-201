---
title: "Task 1: Call API with Postman"
linkTitle: "Task 1: Call API with Postman"
weight: 10
---

|                            |    |  
|----------------------------| ----
| **Goal**                   | Call Juice Shop API with Postman
| **Task**                   | Configure Postman and GET information about Apple Juice from the product page
| **Verify task completion** | The Postman Request should successfully return data for the Apple Juice product

1.  Open postman by opening a new terminal (not bash) and type ```Postman``` at the prompt.  This should start the postman application.
{{% notice warning %}}

If Postman doesn't open, it's likely due to the terminal still using Bash.  To exit bash, simply type ```sh```

{{% /notice %}}

   - When postman opens, select **Continue without an account**
    ![postmanlite](p-light.png)
   - Now select **Open Lightweight API Client**
    ![postmanlite2](p-light2.png)

2.  Now, let's make an HTTP GET API call to search for Apple Juice.  Use the following URL, ensuring you replace your studentID.

    ```sh
    https://<studentID>.fwebtraincse.com/rest/products/search?q=Apple
    ```

{{% notice warning %}}
If the first call fails, due to a certificate error.  In the response section, you will need to scroll down and select "Disable SSL Verification".
![postman ssl disable](p-dis.png)
{{%/notice%}}

3. Now the Call should go through an you should see a status 200 and returned data.


    ![postman success](p-success.png)

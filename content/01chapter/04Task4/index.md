---
title: "Task 4: Check Juice Shop"
linkTitle: "Task 4: Check Juice Shop"
weight: 40
---

|                            |    |  
|----------------------------| ----
| **Goal**                   | Verify that Juice Shop is working
| **Task**                   | Navigate to the public IP associated with Ubuntu
| **Verify task completion** | You should see the Juice Shop Home Page



### Start Kali RDP

From your Terraform Outputs in Task 2 you should have seen gotten the Public IP address of Ubuntu.

- By default, Juice shop listens on port 3000.  In your favorite browser, type ```http://<ubuntu-ip>:3000``` 
- You should see a screen like below:

{{% notice warning %}} Depending on your browser, you might have to acknowledge a warning about the lack of encryption.  {{% /notice %}}

{{< figure src="js-initial.png" alt="Juice Shop" >}}

- You can now proceed to the next module

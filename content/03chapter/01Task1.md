---
title: "Task 1 - Perform a simple SQL injection attack"
linkTitle: "Task 1 - Perform a simple SQL injection attack"
weight: 10
---


|                            |    |  
|----------------------------| ----
| **Goal**                   | Perform SQL injection attack
| **Task**                   | Modify the hosts file on Kali and then perform a simple browser based attack
| **Verify task completion** | You should see that the SQLi attack is accepted by Juice Shop.

1. Log into Kali linux: https://{{Kali IP}}/vnc.html

2. DNS modification takes too long to propagate for this lab, so we will add a host entry on Kali.  Open the terminal emulator on Kali Linux screen. At the prompt, type:

```sh
bash
sudo nano /etc/hosts
```

3. Add a host entry in the format ```<ip_address> <studentId>.fwebtraincse.com``` at the bottom of the file.  
   - <ip_address> will come from the IP addresses provided by FortiWeb Cloud during application setup Chapter 2 Task 1 "Step 5" of the onboarding   
   - <studentId> comes from your FortiWeb Cloud account ID (used to create the FortiWeb application)
   - Once this is complete, type **<kbd>ctrl</kbd>+o** followed by <kbd>enter</kbd> and then **<kbd>ctrl</kbd>+x**

![Hosts](hosts.png)

4. Now let’s Navigate to the Firefox browser (located at the top of Kali desktop) and enter our FortiWeb Cloud Protected Juice Shop URL into the navigation bar ```https://number.fwebtraincse.com```.  Accept warnings and proceed to the application

{{% notice info %}}If we had modified the DNS record at the beginning of this lab, FortiWeb would have pulled a valid SSL certificate from Let's Encrypt.{{% /notice %}}

![stud-home](studhome.png)

5. Let’s perform a very simple SQLi attack. To perform a SQLi attack append ```?name=' OR 'x'='x``` to your URL.  Be sure that you use **YOUR NUMBER**.

   - For example (be sure to use your studentId)
     - ```https://669.fwebtraincse.com/?name=' OR 'x'='x```


{{% notice info %}}
The attack will go through and you will see the Juice Shop Home page
{{% /notice %}}


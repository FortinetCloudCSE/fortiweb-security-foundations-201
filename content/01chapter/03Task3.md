---
title: "Task 3 - Verify Lab Environment"
menuTitle: "Task 3 - Verify Lab Environment"
weight: 30
---

Below is a diagram of the Lab environment.

![lab1](diagram.png)

### Check Availability of Juice Shop

Use the public IP of Juice Shop to log in: http://{{Juice Shop IP}}:3000

![Juiceshop Home Page](juice-home.png)

### Login to Kali

1.  Use the Kali public IP to log in: https://{{Kali IP}}/vnc.html

Accept certificate errors and proceed.  When prompted, click **Connect**.  This will take you to the home screen of Kali

![Kali Home Page](kali-home.png)

2.  In order to copy/paste into Kali, we will need to click on the tab at the left hand side of the screen.

![cp-tab](cp-tab-kali.png)

3.  This will open the tab revealing a couple of options.  Select the clipboard icon and paste your text into the box.  Once the text is in the box, you can right click on the desktop and select Paste Selection to paste in the text.  When done, you can click on the arrow to hide the clipboard.

![paste-kali](paste-kali.png)

{{% notice info %}}We are going to make a small change to Kali in order to prepare for later steps. Please open a terminal by clicking on the Black Icon at the bottom of the screen.{{% /notice %}}

 4.  Enter the following:

```sh
bash
echo nameserver 8.8.8.8 >> /etc/resolv.conf
apt-get install nano
wget https://dl.pstmn.io/download/latest/linux_64 -O /tmp/linux_64 && tar xvzf /tmp/linux_64 -C /tmp/ && sudo mv /tmp/Postman /opt/ && sudo ln -s /opt/Postman/app/Postman /usr/local/bin/Postman
```

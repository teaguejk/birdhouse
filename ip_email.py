import smtplib, email, ssl
import datetime
import subprocess
import socket
from email import encoders
from email.mime.base import MIMEBase
from email.mime.text import MIMEText

# ip = subprocess.check_output(["hostname", "-I"])
ip = socket.gethostname()
recipient = 'teaguejk@appstate.edu'
sender    = 'zerodonotreply@gmail.com'
password  = 'birdhouse1!'

subject = 'IP FROM PI ZERO'
body    = str(ip) + ' at ' + str(datetime.datetime.now())

message = MIMEText(body, "plain")
message['Subject'] = subject
message['From']    = sender
message['To']      = recipient

text = message.as_string()
context = ssl.create_default_context()
with smtplib.SMTP_SSL('smtp.gmail.com', 465, context=context) as server:
	server.login(sender, password)
	server.sendmail(sender, recipient, text)
"""
#
# Birdhouse Camera Project
# Jaracah Teague
# Capstone Project
# Appalachian State University
#
"""

# import RPi.GPIO as GPIO
# from picamera import PiCamera
import time
import datetime
import smtplib, email, ssl
import cv2 as cv
import imutils
import argparse
from imutils.video import VideoStream
from email import encoders
from email.mime.base import MIMEBase
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText

# ssl._create_default_https_context = ssl._create_unverified_context
 
#------------------------------------------------------------------------------------------
#Function defs

def start_video():
	"""
	Function to start capturing video if motion is detected at the camera

	"""
	print ('Motion has been detected at camera\n')
	timestamp = time.strftime('%m-%d-%y-%H-%M-%S')  


def send_email(password):
	"""
	Function to send an email if given conditions are met,
	this email will contain an image attachment
	Sends to a list of users
	
	"""

	# set up the message
	sender_email	= 'zeroDoNotReply@gmail.com'
	password		= password

	recepient_emails = [
		'teaguejk@appstate.edu',
		'erinbrzezin@gmail.com'
	]
	
	timestamp 	= time.strftime('%m-%d-%y %H:%M:%S')
	subject 	= 'Automated Email: Motion Detected in Birdhouse ' + timestamp
	body 		= 'Motion Detected in Birdhouse'
	message 	= MIMEMultipart()

	message["From"] 	= sender_email
	message["To"] 		= recepient_emails[0]
	message["Subject"] 	= subject

	message.attach(MIMEText(body, "plain"))
	
	# add the attachment (image) to the message
	filename = 'IMG.jpg'
	with open(filename, "rb") as attachment:
		part = MIMEBase("application", "octet-stream")
		part.set_payload(attachment.read())

	# encoding
	encoders.encode_base64(part)

	part.add_header(
	"Content-Disposition",
	f"attachment; filename= {filename}",
	)

	message.attach(part)
	
	# sending the message
	text = message.as_string()
	context = ssl.create_default_context()
	with smtplib.SMTP_SSL("smtp.gmail.com", 465, context=context) as server:
		server.login(sender_email, password)
		server.sendmail(sender_email, recepient_emails[0], text)


#------------------------------------------------------------------------------------------
# Main section
def main():
	# create the camera object
    # camera = PiCamera()
    # camera.start_preview()
    # camera.capture('./img.jpg')
    # time.sleep(2)
    # camera.stop_preview()

	# while loop
	# send email
	password = input("Please Enter Password\n")
	send_email(password)


if __name__ == "__main__":
	main()



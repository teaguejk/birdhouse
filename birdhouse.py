"""
#
# Birdhouse Camera Project
# Jaracah Teague
# Capstone Project
# Appalachian State University
#
"""

import RPi.GPIO as GPIO
from picamera import PiCamera
import time
import datetime
import smtplib

#Function defs

def start_video():
	"""
	Function to start capturing video if motion is detected at the camera

	"""
	print ("Motion has been detected at camera\n")
	timestamp = time.strftime("%m-%d-%y-%H-%M-%S")  


def send_email():
	"""
	Function to send an email if given conditions are met
	Sends to a list of users
	
	"""

# Main section

# create the camera object
camera = PiCamera()

# while loop

"""
#
# Birdhouse Camera Project
# Jaracah Teague
# Capstone Project
# Appalachian State University
#
"""

# img/video files stored in ./assets
# pi camera related things may be temporarily commented out 

# py imports
# import pandas as pd
# import numpy as np
import time
import datetime
import smtplib, email, ssl
import cv2
import argparse
import imutils
from email import encoders
from email.mime.base import MIMEBase
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText

# pi imports 
# import RPi.GPIO as GPIO
# from picamera import PiCamera
# from picamera.array import PiRGBArray

# Settings
camera_delay = 2.5
resolution = [640, 480]
fps = 16
min_area = 5000

#------------------------------------------------------------------------------------------
# Notes
"""
# ssl._create_default_https_context = ssl._create_unverified_context

# save image as an array 
raw_capture = PiRGBArray(camera, size=tuple(resolution))

# open cv loop
for f in camera.capture_continuous(raw_capture, format="bgr", use_video_port=True ):
    # get the numpy array that represents the image
    frame = f.array
    timestamp = datetime.datetime.now()
    state = 'Unoccupied'

# create the camera object and capture
    # camera = PiCamera()
    # camera.start_preview()
    # camera.capture('./img.jpg')
    # time.sleep(2)
    # camera.stop_preview()

"""

#------------------------------------------------------------------------------------------
# Function defs

def start_video():
    """
    Function to start capturing video 

    """
    # save a timestamp
    # timestamp = time.strftime('%m-%d-%y-%H-%M-%S')  
    timestamp = datetime.datetime.now()

    # camera setup
    # camera = PiCamera()
    # camera.resolution = tuple(resolution)
    # camera.framerate = fps

    # print message
    print("[MSG] starting camera...")

    # delay start
    time.sleep(camera_delay)

    # start video that will stop on pressing q or a time limit passes

    

def capture_image():
    """
    Function that captures a single image when called using pi camera

    """

    # save a timestamp
    # timestamp = time.strftime('%m-%d-%y-%H-%M-%S')  
    timestamp = datetime.datetime.now()
    filename = f'./assets/{timestamp}_bird.jpg'

    # create camera object
    # camera = PiCamera()

    # # start camera, capture, delay, close
    # camera.start_preview()
    # camera.capture(filename)
    # time.sleep(2)
    # camera.stop_preview()
    
    return filename


def send_email(password, filename):
    """
    Function to send an email if given conditions are met,
    this email will contain an image attachment
    Sends to a list of users
    
    """

    # set up the message
    sender_email    = 'zeroDoNotReply@gmail.com'
    password        = password

    # recipient email list
    # will eventually be retrieved from website
    # by reading a csv file into a pandas dataframe
    recepient_emails = [
        'teaguejk@appstate.edu',
        'erinbrzezin@gmail.com'
    ]
    
    timestamp   = time.strftime('%m-%d-%y %H:%M:%S')
    subject     = 'Automated Email: Motion Detected in Birdhouse ' + timestamp
    body        = 'Motion Detected in Birdhouse'
    message     = MIMEMultipart()

    message["From"]     = sender_email
    message["To"]       = recepient_emails[0]
    message["Subject"]  = subject

    message.attach(MIMEText(body, "plain"))
    
    # add the attachment (image) to the message
    # filename = 'IMG.jpg'
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
# Main routine
def main():
    # get password on start
    password = input("Please Enter Password\n")

    # variables needed for the main loop
    exit = False
    motion_detected = False

    # while not exit:
        # constantly check for motion to be detected until an exit command is entered or ctrl-c

        # if motion_detected:
            # if motion detected, capture image and pass the filename to send_email
            # filename = capture_image()
            # send_email(password, filename

    pass
    
if __name__ == "__main__":
    main()



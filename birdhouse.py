"""
#
# Birdhouse Camera Project
# Jaracah Teague
# Capstone Project
# Appalachian State University
#
# img/video files stored in ./assets and named in this format: {timestamp}_bird.jpg/mov
# 
# assumes that the student2.cs.appstate.edu server is hosting the site 
# located at usr/local/apache2/htdocs/u/teaguejk/birdhouse.site
# 
# pi camera related things may be temporarily commented out 
#
# emails are sent from 'zeroDoNotReply@gmail.com' but can be changed
#
# passwords are entered on runtime for security
#
"""

#==========================================================================================
# Setup

#------------------------------------------------------------------------------------------
# py imports
import pandas as pd
import numpy as np

import time
import datetime
import argparse

import cv2
import imutils

import smtplib, email, ssl
from email import encoders
from email.mime.base import MIMEBase
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText

from fabric.connection import Connection

# pi imports 
# Comment out when not on PI
import RPi.GPIO as GPIO
from picamera import PiCamera
from picamera.array import PiRGBArray

#------------------------------------------------------------------------------------------
# Settings

# GPIO Settings
MOTION_PIN = 4

# Camera Settings
camera_delay = 1
resolution = [640, 480]
fps = 16
min_area = 5000

# Web Server Setup
user = 'teaguejk'
host = 'student2.cs.appstate.edu'
path = '/usr/local/apache2/htdocs/u/teaguejk/birdhouse.site'
# spass = ""
# Mailing list is stored in a csv file at /usr/local/apache2/htdocs/u/teaguejk/birdhouse.site/mailing_list.csv

# Sender Email
sender_email = 'zeroDoNotReply@gmail.com'
# epass = ""
#------------------------------------------------------------------------------------------

#==========================================================================================
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

# create the camera object and capture
    # camera = PiCamera()
    # camera.start_preview()
    # camera.capture('./img.jpg')
    # time.sleep(2)
    # camera.stop_preview()

# server.sendmail(sender_email, recepient_emails[0], text) # send an email



"""

#==========================================================================================
# Function defs

def send_email(epass, spass, filename):
    """
    Function to send an email if given conditions are met,
    this email will contain an image attachment
    Sends to a list of users
    
    """
    # recipient email list
    # will eventually be retrieved from website
    # by reading a csv file into a pandas dataframe
    with Connection(host, user, connect_kwargs={'password': spass, 'allow_agent': False}) as c, c.sftp() as sftp,   \
         sftp.open(path + '/mailing_list.csv') as file:
            mailing_list = pd.read_csv(file, skip_blank_lines=True, header=None).T
    # drop NaN or None type entries
    mailing_list = mailing_list.dropna() 

    # message setup

    timestamp   = time.strftime('%m-%d-%y %H:%M:%S')
    subject     = 'Automated Email: Motion Detected in Birdhouse ' + timestamp
    body        = 'Motion Detected in Birdhouse'
    message     = MIMEMultipart()

    message["From"]     = sender_email
    message["To"]       = sender_email
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
    # if certificates cannot be verified, use this context instead of the default
    # ssl.SSLContext.verify_mode = ssl.VerifyMode.CERT_NONE
    # context = ssl.SSLContext(ssl.PROTOCOL_TLS)
    context = ssl.create_default_context()
    with smtplib.SMTP_SSL("smtp.gmail.com", 465, context=context) as server:
        server.login(sender_email, epass)
        # loop through dataframe and send email to each row (address) 
        for row in mailing_list.iterrows(): 
                server.sendmail(sender_email, row[1].values, text)

#------------------------------------------------------------------------------------------
# Works only on PI

def capture_video():
    """
    Function to capture a 60 second video 
    Comment out when not on PI


    """
    # timestamp
    # # timestamp = time.strftime('%m-%d-%y-%H-%M-%S')  
    # timestamp = datetime.datetime.now()

    # filename
    # filename = f'./assets/{timestamp}_bird.mov'

    # # camera setup
    # camera = PiCamera()
    # camera.resolution = tuple(resolution)
    # camera.framerate = fps

    # # print message
    # print("[MSG] Starting Camera...\n")
    # # delay start
    # time.sleep(camera_delay)
    
    # # record
    # print("[MSG] Starting Recording...\n")
    # camera.start_recording(filename)
    # camera.wait_recording(60)
    # camera.stop_recording()
    # printf("[MSG] Captured Video...\n")

    pass

def capture_image():
    """
    Function that captures a single image when called using pi camera
    Comment out when not on PI

    """

    # timestamp
    # timestamp = time.strftime('%m-%d-%y-%H-%M-%S')  
    timestamp = datetime.datetime.now()

    # filename
    filename = f'./assets/{timestamp}_bird.jpg'
    # filename = f'./assets/IMG.jpg'

    # start camera, capture, delay, close
    print("[MSG] Starting Camera...\n")
    time.sleep(camera_delay)

    camera = PiCamera()

    camera.start_preview()
    camera.capture(filename)
    time.sleep(2)
    camera.stop_preview()
    print("[MSG] Captured Image...\n")

    # upload the new image into the directory
    with Connection(host, user, connect_kwargs={'password': spass, 'allow_agent': False}) as c:  
         c.put(filename, path + "/assets/IMG.jpg")

    return filename

    # pass

#------------------------------------------------------------------------------------------

#==========================================================================================
# Main routine

def main():
    # get passwords on start
    print("For Email: "  + sender_email)
    epass = input("Email Password:\n")
    print("For server: " + user + "@" + host)
    spass = input("Server Password:\n")

    # motion detection from GPIO pin (MOTION_PIN)
    GPIO.setmode(GPIO.BCM)
    GPIO.setup(MOTION_PIN, GPIO.IN)

    # when motion is detected, capture_image (upload new IMG to server) -> send_email
    # exits when CTRL-C is typed
    try:
        print("[MSG] Starting Motion Detection\n")
        print('[MSG] Type CTRL-C to exit\n')
        time.sleep(2)
        print('[MSG] Active\n')
        while True:
            if GPIO.input(MOTION_PIN):
                print("[MSG] Motion Detected\n")
                # capture image -> upload to server, overwriting old IMG -> save filename
                filename = capture_image()
                # pass in the email password, server password, and new filename
                send_email(epass, spass, filename)
                time.sleep(1)
    except KeyboardInterrupt:
        print("[MSG] Closing\n")
        GPIO.cleanup()
        time.sleep(2)
        print("[MSG] Inactive\n")


    # For testing
    # filename = './assets/IMG.jpg'
    # filename = capture_image()
    # send_email(epass, spass, filename)


    pass
    
if __name__ == "__main__":
    main()

#==========================================================================================


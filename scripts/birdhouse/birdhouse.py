"""
#
# Birdhouse Camera Project
# Jaracah Teague
# Appalachian State University
#
"""


import time
import datetime
import smtplib
import ssl
import pandas as pd
from email import encoders
from email.mime.base import MIMEBase
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from fabric.connection import Connection

import RPi.GPIO as GPIO
from picamera import PiCamera


# GPIO Settings
MOTION_PIN = 4

# Camera Settings
camera_delay = 1
resolution = [640, 480]
fps = 16
min_area = 5000

# Web Server Setup
user = 'teaguejk'
path = ''
# Mailing list is stored in a csv file at <path>/mailing_list.csv

# Sender Email
sender_email = 'zeroDoNotReply@gmail.com'
# epass = ""


def send_email(epass, filename):
    """
    Function to send an email if given conditions are met,
    this email will contain an image attachment
    Sends to a list of users
    
    """
    # recipient email list
    # will eventually be retrieved from website
    # by reading a csv file into a pandas dataframe
    with open('/mailing_list.csv') as file:
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


def capture_image():
    """
    Function that captures a single image when called using pi camera
    Comment out when not on PI

    """

    timestamp = datetime.datetime.now()
    filename = f'./assets/{timestamp}_bird.jpg'

    # start camera, capture, delay, close
    print("[MSG] Starting Camera...\n")
    time.sleep(camera_delay)

    camera = PiCamera()

    camera.start_preview()
    camera.capture(filename)
    time.sleep(2)
    camera.stop_preview()

    print("[MSG] Captured Image...\n")
    print("[MSG] Closing Camera...\n")
    camera.close()

    return filename

    # pass


def main():
    # get passwords on start
    print("Email: "  + sender_email)
    epass = input("Email Password:\n")

    # motion detection from GPIO pin (MOTION_PIN)
    GPIO.setmode(GPIO.BCM)
    GPIO.setup(MOTION_PIN, GPIO.IN)

    # exits when CTRL-C is typed
    try:
        print("[MSG] Starting Motion Detection\n")
        print('[MSG] Type CTRL-C to exit\n')
        time.sleep(2)
        print('[MSG] Active\n')
        while True:
            if GPIO.input(MOTION_PIN):
                print("[MSG] Motion Detected\n")
                filename = capture_image()

                send_email(epass, filename)

                print('Sleeping...\n')
                time.sleep(1)
                print('[MSG] Active\n')
    except KeyboardInterrupt:
        print("[MSG] Closing\n")
        GPIO.cleanup()
        time.sleep(120)
        print("[MSG] Inactive\n")

if __name__ == "__main__":
    main()

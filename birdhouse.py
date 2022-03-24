"""
#
# Birdhouse Camera Project
# Jaracah Teague
# Capstone Project
# Appalachian State University
#
"""
# img/video files stored in ./assets

# py
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

# pi
# import RPi.GPIO as GPIO
from picamera import PiCamera
from picamera.array import PiRGBArray


# misc
# ssl._create_default_https_context = ssl._create_unverified_context

# Settings
camera_delay = 2.5
resolution = [640, 480]
fps = 16
min_area = 5000

#------------------------------------------------------------------------------------------
# Function defs

def start_video():
    """
    Function to start capturing video 

    """
    # print ('Motion has been detected at camera\n')
    timestamp = time.strftime('%m-%d-%y-%H-%M-%S')

    # camera setup
    camera = PiCamera()
    camera.resolution = tuple(resolution)
    camera.framerate = fps

    print("[MSG] starting camera...")
    time.sleep(camera_delay)

    rawCapture = PiRGBArray(camera, size=tuple(resolution))

    firstFrame = None

    for f in camera.capture_continuous(rawCapture, format="bgr", use_video_port=True):
        timestamp = datetime.datetime.now()
        rawCapture.truncate(0)
        rawCapture.seek(0)

        # put frame into array and resize
        frame = f.array
        frame = imutils.resize(frame, width=500)

        # convert to grayscale
        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        gray = cv2.GaussianBlur(gray, (21, 21), 0)        
        
        if firstFrame is None:
            firstFrame = gray
            continue

        # frameDelta - background_model - currentFrame
	    frameDelta = cv2.absdiff(firstFrame, gray)

        # find and draw contours

        # display any text on the frame
        ts = timestamp.strftime("%A %d %B %Y %I:%M:%S%p")
        cv2.putText(frame, ts, (10, frame.shape[0] - 10), cv2.FONT_HERSHEY_SIMPLEX,
            0.35, (0, 0, 255), 1)

        

    # for frame in camera.capture_continuous(rawCapture, format="bgr", use_video_port=True):
    #     image = frame.array
    #     result = predictor.detect(image)
    #     for obj in result:
    #         logger.info('coordinates: {} {}. class: "{}". confidence: {:.2f}'.format(obj[0], obj[1], obj[3], obj[2]))
    #         cv2.rectangle(image, obj[0], obj[1], (0, 255, 0), 2)
    #         cv2.putText(image, '{}: {:.2f}'.format(obj[3], obj[2]), (obj[0][0], obj[0][1] - 5), cv2.FONT_HERSHEY_PLAIN, 1, (0, 255, 0), 2)

    #     # show the frame
    #     cv2.imshow("Stream", image)
    #     key = cv2.waitKey(1) & 0xFF


    # camera.start_preview()
    # camera.capture('./assets/img.jpg')
    # time.sleep(2)
    # camera.stop_preview()


def send_email(password):
    """
    Function to send an email if given conditions are met,
    this email will contain an image attachment
    Sends to a list of users
    
    """

    # set up the message
    sender_email    = 'zeroDoNotReply@gmail.com'
    password        = password

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
    # camera.capture('./assets/img.jpg')
    # time.sleep(2)
    # camera.stop_preview()
    
    # while loop
    # send email
    # password = input("Please Enter Password\n")
    # send_email(password)
    
    start_video()

if __name__ == "__main__":
    main()



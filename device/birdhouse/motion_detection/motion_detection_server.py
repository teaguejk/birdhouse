from flask import Flask, render_template, Response
import cv2
import time
import datetime
import numpy as np
from picamera2 import Picamera2
import os

app = Flask(__name__)

capture_dir = "capture"

picam2 = Picamera2()
picam2.configure(picam2.create_preview_configuration(
    main={"format": 'XRGB8888', "size": (640, 480)}))
picam2.start()

time.sleep(2)

avg_frame = None

os.makedirs(capture_dir, exist_ok=True)

def gen_frames():
    global avg_frame
    while True:
        frame = picam2.capture_array()

        # greyscale and blur to reduce noise
        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        gray = cv2.GaussianBlur(gray, (21, 21), 0)

        # keep average frame
        if avg_frame is None:
            avg_frame = gray.astype("float")
            continue

        # accumulate weight average
        cv2.accumulateWeighted(gray, avg_frame, 0.5)

        # diff between current and average frame
        frame_delta = cv2.absdiff(gray, cv2.convertScaleAbs(avg_frame))

        # apply threshold to delta and dilate to fill in holes
        thresh = cv2.threshold(frame_delta, 25, 255, cv2.THRESH_BINARY)[1]
        thresh = cv2.dilate(thresh, None, iterations=2)

        # find contours
        contours, _ = cv2.findContours(
            thresh.copy(), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)

        motion_detected = False

        for contour in contours:
            if cv2.contourArea(contour) < 500:
                continue

            motion_detected = True

            # bounding box
            (x, y, w, h) = cv2.boundingRect(contour)
            cv2.rectangle(frame, (x, y), (x + w, y + h), (0, 255, 0), 2)

        if motion_detected:
            cv2.putText(frame, "Motion Detected", (10, 20),
                        cv2.FONT_HERSHEY_SIMPLEX, 0.6, (0, 0, 255), 2)
            
            # save the image
            timestamp = datetime.datetime.now().strftime("%Y%m%d-%H%M%S.%f")
            filename = f"motion_{timestamp}.jpg"
            save_path = os.path.join(capture_dir, filename)

            cv2.imwrite(save_path, frame)

        # Encode the frame in JPEG format
        ret, buffer = cv2.imencode('.jpg', frame)
        frame = buffer.tobytes()

        # Yield the frame in byte format
        yield (b'--frame\r\n' b'Content-Type: image/jpeg\r\n\r\n' + frame + b'\r\n')

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/video_feed')
def video_feed():
    # Return the response generated along with the specific media type (mime type)
    return Response(gen_frames(),
                    mimetype='multipart/x-mixed-replace; boundary=frame')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, threaded=False)

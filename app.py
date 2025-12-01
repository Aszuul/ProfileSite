from flask import Flask, render_template, url_for, request, jsonify, redirect, flash
import requests
import os
from flask_mail import Mail, Message

CAPTCHA_SECRET_KEY=os.environ['CAPTCHA_SECRET_KEY']
CAPTCHA_SITE_KEY=os.environ['CAPTCHA_SITE_KEY']
SMTP_KEY=os.environ['SMTP_KEY']

app = Flask(__name__)
app.config['SECRET_KEY'] = os.environ.get('SECRET_KEY', 'very secret')
app.config['MAIL_SERVER'] = 'smtp.gmail.com'
app.config['MAIL_PORT'] = 587
app.config['MAIL_USE_TLS'] = True
app.config['MAIL_USE_SSL'] = False
app.config['MAIL_USERNAME'] = 'scott.abernathy20@gmail.com'
app.config['MAIL_PASSWORD'] = SMTP_KEY
app.config['MAIL_DEFAULT_SENDER'] = 'noreply@scott-abernathy.com'
mail = Mail(app)
FLASK_DEBUG = os.environ.get('FLASK_ENV', True)

def verify_recaptcha(token, action):
    """Verifies the reCAPTCHA token with Google."""
    RECAPTCHA_VERIFY_URL = 'https://www.google.com/recaptcha/api/siteverify'
    
    response = requests.post(RECAPTCHA_VERIFY_URL, data={
        'secret': CAPTCHA_SECRET_KEY,
        'response': token,
        'remoteip': request.remote_addr # Optional but recommended
    })
    
    result = response.json()
    
    # Check if verification succeeded and the action name matches
    if not result.get('success'):
        return False, "Verification failed (Google response failed)."

    # Check the score and the action name
    MIN_SCORE = 0.5 # Set your acceptable threshold (0.5 is default)
    if result.get('score', 0.0) < MIN_SCORE:
        return False, f"Score too low: {result.get('score')} < {MIN_SCORE}"
        
    if result.get('action') != action:
        # This prevents an attacker from using a token from a different action (e.g., 'login')
        return False, f"Action mismatch: expected '{action}', got '{result.get('action')}'"
        
    return True, "Success"

def send_mail(name, message, email):
    msg = Message(
        f'Contact Me - {name}: {email}',
        recipients=['scott.abernathy20@gmail.com'],
        body=f'{message}',
    )
    mail.send(msg)
    return 'Email sent successfully!'

@app.route('/')
def home():
    return render_template('index.html', CAPTCHA_SITE_KEY=CAPTCHA_SITE_KEY)

@app.route('/about')
def about():
    return render_template('about.html', CAPTCHA_SITE_KEY=CAPTCHA_SITE_KEY)

@app.route('/projects')
def projects():
    return render_template('projects.html', CAPTCHA_SITE_KEY=CAPTCHA_SITE_KEY)

@app.route('/sendmail', methods=['POST'])
def sendmail():
    # validate recaptcha
    token = request.form.get('g-recaptcha-response')
    if not token: 
        return jsonify({
            'success': False,
            'message': 'reCAPTCHA token is missing'
        }), 400
    valid, res = verify_recaptcha(token, 'contact_submit')
    if not valid:
        flash('There was an issue with your submission.', 'danger')
        return redirect(url_for())
    name = request.form.get('name')
    email = request.form.get('email')
    message = request.form.get('message')
    send_mail(name, message, email)
    flash('Email submitted successfully!', 'success')
    if request.referrer:
        return redirect(request.referrer)
    else:
        return redirect(url_for('home'))

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=80, debug=FLASK_DEBUG)
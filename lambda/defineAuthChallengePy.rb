import json

def lambda_handler(event, context):
    if event['request']['session']:
        previous = event['request']['session'][-1]
        if previous['challengeResult']:
            event['response']['issueTokens'] = True
            event['response']['failAuthentication'] = False
        else:
            event['response']['issueTokens'] = False
            event['response']['failAuthentication'] = True
    else:
        event['response']['issueTokens'] = False
        event['response']['failAuthentication'] = False
        event['response']['challengeName'] = 'CUSTOM_CHALLENGE'
    return event

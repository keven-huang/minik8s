# function3 - sub .py
def Function3(event: dict, context: dict)->dict:
    result = context.get('result')

    if result is not None:
        try:
            result = result - 10
            return {'result': result}
        except Exception as e:
            return {'error': str(e)}

    return {'error': 'Invalid input'}

def chunks(list, count):
    """Yield successive n-sized chunks from l."""
    for i in xrange(0, len(list), count):
        yield list[i:i + count]
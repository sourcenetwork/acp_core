package zanzi

type loggerWrapper struct {
}

func (l *loggerWrapper) Infof(fmt string, args ...any) {
}

func (l *loggerWrapper) Debugf(fmt string, args ...any) {
}

func (l *loggerWrapper) Errorf(fmt string, args ...any) {
}

func (l *loggerWrapper) Warnf(fmt string, args ...any) {
}

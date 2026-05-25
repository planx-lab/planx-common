package metrics

import "testing"

func TestNoopCounter_Inc(t *testing.T) {
	c := NoopCounter{}
	c.Inc()
}

func TestNoopCounter_Add(t *testing.T) {
	c := NoopCounter{}
	c.Add(42.0)
}

func TestNoopCounter_Interface(t *testing.T) {
	var _ Counter = NoopCounter{}
}

func TestNoopGauge_Set(t *testing.T) {
	g := NoopGauge{}
	g.Set(100.0)
}

func TestNoopGauge_Inc(t *testing.T) {
	g := NoopGauge{}
	g.Inc()
}

func TestNoopGauge_Dec(t *testing.T) {
	g := NoopGauge{}
	g.Dec()
}

func TestNoopGauge_Add(t *testing.T) {
	g := NoopGauge{}
	g.Add(10.0)
}

func TestNoopGauge_Sub(t *testing.T) {
	g := NoopGauge{}
	g.Sub(5.0)
}

func TestNoopGauge_Interface(t *testing.T) {
	var _ Gauge = NoopGauge{}
}

func TestNoopHistogram_Observe(t *testing.T) {
	h := NoopHistogram{}
	h.Observe(1.23)
}

func TestNoopHistogram_Interface(t *testing.T) {
	var _ Histogram = NoopHistogram{}
}

func TestNoopProvider_Counter(t *testing.T) {
	p := NoopProvider{}
	c := p.Counter("test_counter", map[string]string{"k": "v"})
	if c == nil {
		t.Fatal("Counter returned nil")
	}
}

func TestNoopProvider_Gauge(t *testing.T) {
	p := NoopProvider{}
	g := p.Gauge("test_gauge", map[string]string{"k": "v"})
	if g == nil {
		t.Fatal("Gauge returned nil")
	}
}

func TestNoopProvider_Histogram(t *testing.T) {
	p := NoopProvider{}
	h := p.Histogram("test_hist", map[string]string{"k": "v"})
	if h == nil {
		t.Fatal("Histogram returned nil")
	}
}

func TestNoopProvider_Interface(t *testing.T) {
	var _ Provider = NoopProvider{}
}

func TestNoopProvider_EmptyLabels(t *testing.T) {
	p := NoopProvider{}
	c := p.Counter("c", nil)
	if c == nil {
		t.Fatal("Counter with nil labels returned nil")
	}
}

func TestNoopProvider_EmptyName(t *testing.T) {
	p := NoopProvider{}
	g := p.Gauge("", map[string]string{})
	if g == nil {
		t.Fatal("Gauge with empty name returned nil")
	}
}

func TestNoopCounter_DoesNotPanic(t *testing.T) {
	c := NoopCounter{}
	c.Add(-1.0)
	c.Add(0)
	c.Inc()
}

func TestNoopGauge_DoesNotPanic(t *testing.T) {
	g := NoopGauge{}
	g.Set(0)
	g.Add(-100)
	g.Sub(-100)
	g.Inc()
	g.Dec()
}

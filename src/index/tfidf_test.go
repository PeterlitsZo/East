package index

import (
    "testing"
    "math"
)

func simMap(a map[string]float64, b map[string]float64) bool {
    for word, _ := range a {
        if math.Abs(b[word] - a[word]) > 0.1 {
            return false
        }
    }
    for word, _ := range b {
        if math.Abs(b[word] - a[word]) > 0.1 {
            return false
        }
    }
    return true
}

func TestGetTf_1 (t *testing.T) {
    tf := GetTf("balabalabala")
    expended := map[string]float64{"balabalabala": 1}
    if !simMap(tf, expended) {
        t.Errorf("I get %#v, but I want %#v", tf, expended)
    }
}

func TestGetTf_2 (t *testing.T) {
    tf := GetTf("balabalabala a a a a a")
    expended := map[string]float64{
        "a": 2.61,
        "balabalabala": 1,
    }
    if !simMap(tf, expended) {
        t.Errorf("I get %#v, but I want %#v", tf, expended)
    }
}

func TestGetTf_3 (t *testing.T) {
    tf := GetTf("who is peter? maybe you want to ask, but the true thing is peter is a person")
    expended := map[string]float64{
        "a": 1,
        "ask": 1,
        "but": 1,
        "is": 2.10,
        "maybe": 1,
        "person": 1,
        "peter": 1.69,
        "the": 1,
        "thing": 1,
        "to": 1,
        "true": 1,
        "want": 1,
        "who": 1,
        "you": 1,
    }
    if !simMap(tf, expended) {
        t.Errorf("I get %#v, but I want %#v", tf, expended)
    }
}


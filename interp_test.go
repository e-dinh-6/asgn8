package main

import "testing"

//example test

func NumCtoV(t *testing.T){
    got1, got2 := interp(NumC{12}, topEnv)
    want1 := NumV{12} 

    if got1 != want1 {
        t.Errorf("got %v, wanted %v", got1, want1)
    }

    if got2 !=  nil {
        t.Errorf("got %v, wanted %v", got2, nil)
    }
}

func StrCtoV(t *testing.T){
    got1, got2 := interp(StrC{"hi"}, topEnv)
    want := StrV{"hi"}

    if got1 != want {
        t.Errorf("got %v, wanted %v", got1, want)
    }

    if got2 != nil {
        t.Errorf("got %v, wanted %v", got2, nil)
    }
}

func IfC_Lookup(t *testing.T){
    got1, got2 := interp(IdC{"a"}, topEnv)
    want := NumC {12}

    if got1 != want {
        t.Errorf("got %v, wanted %v", got1, want)
    }
    if got2 != nil {
        t.Errorf("got %v, wanted %v", got2, nil)
    }
}

// func LamCtoLamV(t *testing.T){
//     got := interp((LamC ((IdC 'a')), (AppC '+', ((IdC 'a'), (NumC 2)))), topEnv)
//     want := (CloV ['a'],  )

//     if got != want {
//         t.Errorf("got %q, wanted %q", got, want)
//     }
// }

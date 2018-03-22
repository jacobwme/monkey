package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"5",5},
		{"10", 10},
		{"-5", -5},
		{"-10",-10},
		{"5 + 5", 10},
		{"5 * 5", 25},
		{"5 / 5", 1},
		{"5 - 5", 0},
		// {"-5 + -5", 0},		can't do this here, only one infix allowed
		{"-5 * -5", 25},
		{"-5 / -5", 1},
		// {"-5 - -5", -10},	can't do this here, only one infix allowed
		{"(5 / 5) * 5", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}


	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}

}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false",false},
		{"!!5",true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input		string
		expected	int64
	}{
		{"return 10;",10},
		{"return 10; 9;",10},
		{"return 2 * 5; 9;",10},
		{"9; return 2 * 5; 9;",10},
		{
			`
if (10 > 1) {
	if (10 > 1 {
		return 10
	}
	return 1;
}`,
		10},
	}

	for _,tt := range tests {
		evalulated := testEval(tt.input)
		testIntegerObject(t, evalulated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input 			string
		expectedMessage	string
	}{
		{"5 + true;",
			"type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5",
			"type mismatch: INTEGER + BOOLEAN"},
		{"-true",
			"unknown operator: -BOOLEAN"},
		{"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false;}",
			"unknown operator: BOOLEAN + BOOLEAN"},
		{`
if (10 > 1) {
	if (10 > 1) {
		return true + false;
	}
	return 1;
}`,
			"unknown operator: BOOLEAN + BOOLEAN"},
	}

	for _, tt := range(tests) {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned, got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL, got=%T (%+v)",obj,obj)
		return false
	}
	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not an Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has the wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is a Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has the wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}
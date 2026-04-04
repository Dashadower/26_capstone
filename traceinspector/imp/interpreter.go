package imp

import "fmt"

type ImpState struct {
	vars         map[string]ImpValues // all exprs are reduced to go values
	return_value ImpValues            // The return value of the current local scope, if exists
}

type ImpInterpreter struct {
	states    []*ImpState // a stack of program states
	functions map[string]*Stmt
}

func (interpreter *ImpInterpreter) get_top_state() *ImpState {
	return interpreter.states[len(interpreter.states)-1]
}

func (interpreter *ImpInterpreter) push_state(state ImpState) {
	interpreter.states = append(interpreter.states, &state)
}

func (interpreter *ImpInterpreter) pop_state() {
	interpreter.states = interpreter.states[:len(interpreter.states)-1]
}

func (interpreter *ImpInterpreter) eval_VarExpr(node VarExpr) ImpValues {
	var_value, var_exists := interpreter.get_top_state().vars[node.name]
	if !var_exists {
		panic("Unknown variable " + node.name)
	}
	return var_value
}

func (interpreter *ImpInterpreter) eval_Expr(node Expr) ImpValues {

}

func (interpreter *ImpInterpreter) eval_Expr_lvalue(lhs Expr, rhs ImpValues) ImpValues {
	lhs_var, lhs_is_var := lhs.(*VarExpr)
	if lhs_is_var {
		_, lhs_exists := interpreter.get_top_state().vars[lhs_var.name]
		if !lhs_exists {
			switch rhs.(type) {
			case *IntVal:
				interpreter.get_top_state().vars[lhs_var.name] = &IntVal{}
			case *BoolVal:
				interpreter.get_top_state().vars[lhs_var.name] = &BoolVal{}
			case *ArrayVal:
				interpreter.get_top_state().vars[lhs_var.name] = &ArrayVal{}
			}
		}
		return interpreter.get_top_state().vars[lhs_var.name]
	} else {
		return interpreter.eval_Expr(lhs)
	}
}

func (interpreter *ImpInterpreter) eval_IntValueExpr(node IntValueExpr) ImpValues {
	return &IntVal{val: node.value}
}

func (interpreter *ImpInterpreter) eval_BoolValueExpr(node BoolValueExpr) ImpValues {
	return &BoolVal{val: node.value}
}

func (interpreter *ImpInterpreter) eval_AddExpr(node AddExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("LHS of addition should be an int value, but got '%s'", node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("RHS of addition should be an int value, but got '%s'", node.rhs))
	}
	return &IntVal{val: lhs_val.val + rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_SubExpr(node SubExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("LHS of subtraction should be an int value, but got '%s'", node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("RHS of subtraction should be an int value, but got '%s'", node.rhs))
	}
	return &IntVal{val: lhs_val.val - rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_MulExpr(node MulExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("LHS of multiplication should be an int value, but got '%s'", node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("RHS of multiplication should be an int value, but got '%s'", node.rhs))
	}
	return &IntVal{val: lhs_val.val * rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_DivExpr(node AddExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("LHS of division should be an int value, but got '%s'", node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("RHS of division should be an int value, but got '%s'", node.rhs))
	}
	return &IntVal{val: lhs_val.val / rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_ArrayIndexExpr(node ArrayIndexExpr) ImpValues {
	index_val, index_is_int := interpreter.eval_Expr(node.index).(*IntVal)
	if !index_is_int {
		panic(fmt.Sprintf("Index of array indexing should be an int value, but got '%s'", node.index))
	}
	base_val, base_is_arrayval := interpreter.eval_Expr(node.base).(*ArrayVal)
	if !base_is_arrayval {
		panic(fmt.Sprintf("Expr %s is not an array", node.base))
	}
	return base_val.val[index_val.val]
}

func (interpreter *ImpInterpreter) eval_EqExpr(node EqExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.lhs)
	rhs_val := interpreter.eval_Expr(node.rhs)
	return &BoolVal{val: lhs_val == rhs_val}
}

func (interpreter *ImpInterpreter) eval_NeqExpr(node NeqExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.lhs)
	rhs_val := interpreter.eval_Expr(node.rhs)
	return &BoolVal{val: lhs_val != rhs_val}
}

func (interpreter *ImpInterpreter) eval_NotExpr(node NotExpr) ImpValues {
	subexpr_val, subexpr_is_bool := interpreter.eval_Expr(node.subexpr).(*BoolVal)
	if !subexpr_is_bool {
		panic(fmt.Sprintf("Subexpr %s of NOT operator should be of type bool", node.subexpr))
	}
	return &BoolVal{val: !subexpr_val.val}
}

func (interpreter *ImpInterpreter) eval_AndExpr(node AndExpr) ImpValues {
	lhs_val, lhs_is_bool := interpreter.eval_Expr(node.lhs).(*BoolVal)
	rhs_val, rhs_is_bool := interpreter.eval_Expr(node.rhs).(*BoolVal)

	if !lhs_is_bool {
		panic(fmt.Sprintf("LHS of AND should be a bool value, but got '%s'", node.lhs))
	}

	if !rhs_is_bool {
		panic(fmt.Sprintf("RHS of AND should be a bool value, but got '%s'", node.rhs))
	}
	return &BoolVal{val: lhs_val.val && rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_OrExpr(node OrExpr) ImpValues {
	lhs_val, lhs_is_bool := interpreter.eval_Expr(node.lhs).(*BoolVal)
	rhs_val, rhs_is_bool := interpreter.eval_Expr(node.rhs).(*BoolVal)

	if !lhs_is_bool {
		panic(fmt.Sprintf("LHS of OR should be a bool value, but got '%s'", node.lhs))
	}

	if !rhs_is_bool {
		panic(fmt.Sprintf("RHS of OR should be a bool value, but got '%s'", node.rhs))
	}
	return &BoolVal{val: lhs_val.val || rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_CallExpr(node CallExpr) ImpValues {
	func_local_state := ImpState{vars: make(map[string]ImpValues)}
	for _, arg := range node.args {
		func_local_state.vars[arg.name] = interpreter.eval_Expr(arg.expr)
	}
	interpreter.push_state(func_local_state)
	interpreter.eval_Stmt(*interpreter.functions[node.func_name])
	return_value := interpreter.get_top_state().return_value
	interpreter.pop_state()
	if return_value == nil {
		return &NoneVal{}
	} else {
		return return_value
	}
}

func (interpreter *ImpInterpreter) eval_Skip(SkipStmt) {}

func (interpreter *ImpInterpreter) eval_AssignStmt(node AssignStmt) {
	rhs_val := interpreter.eval_Expr(node.rhs)
	switch lhs_loc := interpreter.eval_Expr_lvalue(node.lhs, rhs_val).(type) {
	case *IntVal:
		rhs_intval, rhs_is_intval := rhs_val.(*IntVal)
		if !rhs_is_intval {
			panic(fmt.Sprintf("Attempting to assign RHS '%s' of type %T to LHS '%s' of type %T", node.rhs, rhs_val, node.lhs, lhs_loc))
		}
		lhs_loc.val = rhs_intval.val
	case *BoolVal:
		rhs_intval, rhs_is_boolval := rhs_val.(*BoolVal)
		if !rhs_is_boolval {
			panic(fmt.Sprintf("Attempting to assign RHS '%s' of type %T to LHS '%s' of type %T", node.rhs, rhs_val, node.lhs, lhs_loc))
		}
		lhs_loc.val = rhs_intval.val
	case *ArrayVal:
		rhs_intval, rhs_is_arrayval := rhs_val.(*ArrayVal)
		if !rhs_is_arrayval {
			panic(fmt.Sprintf("Attempting to assign RHS '%s' of type %T to LHS '%s' of type %T", node.rhs, rhs_val, node.lhs, lhs_loc))
		}
		lhs_loc.val = rhs_intval.val
	}
}

func (interpreter *ImpInterpreter) eval_Stmt(node Stmt) {

}

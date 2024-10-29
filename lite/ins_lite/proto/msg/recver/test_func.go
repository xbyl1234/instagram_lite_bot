package recver

import (
	"container/list"
	"fmt"
	"github.com/bytedance/sonic/ast"
)

var nameIdx int

type TreeNode struct {
	View   *SubScreen
	Parent *TreeNode
	Child  []*TreeNode
}

func getScreenId(root *SubScreen, useIndex bool) string {
	name := fmt.Sprintf("screen%d", root.Type)
	if useIndex {
		nameIdx++
		name += fmt.Sprintf("_%d_%s_%s", nameIdx, root.GetBaseScreen().WindowId.Value, root.GetBaseScreen().TitleEng.Value)
	} else {
		name += fmt.Sprintf("_%p_%d_%s_%v", root.Value, root.GetBaseScreen().LikeActionResourceId, root.GetBaseScreen().WindowId, root.GetBaseScreen().GetIsLikeResIdChildFlag())
		if root.Type == 1 {
			name += "_" + root.ToSubScreen1().ShowText.Value
		}
	}
	return name
}

func tree2Json(root *TreeNode, useIndex bool) ast.Node {
	jsonObj := ast.NewObject(nil)
	for i := 0; i < len(root.Child); i++ {
		jsonObj.Set(getScreenId(root.Child[i].View, useIndex), tree2Json(root.Child[i], useIndex))
	}
	return jsonObj
}

func Tree2String(root *TreeNode, useIndex bool) string {
	nameIdx = 0
	if root == nil {
		return "root is null"
	}
	jsonObj := ast.NewObject(nil)
	jsonObj.Set(getScreenId(root.View, useIndex), tree2Json(root, useIndex))
	marshal, err := jsonObj.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(marshal)
}

func DumpViewTreeRecvView(root *ScreenReceived) *TreeNode {
	return DumpViewTree(WrapSubScreen(&root.DecodeBody.SubScreen2))
}

func DumpViewTree(root *SubScreen) *TreeNode {
	stack := list.New()
	treeRoot := &TreeNode{View: root}
	stack.PushBack(treeRoot)
	for stack.Len() > 0 {
		cur := stack.Back()
		stack.Remove(cur)
		view := cur.Value.(*TreeNode).View
		if view.Value == nil {
			continue
		}
		allSub := view.GetSubScreenRaw()
		if allSub == nil || allSub.Count() == 0 {
			continue
		}
		for i := 0; i < allSub.Count(); i++ {
			subView := allSub.Get(i)
			node := &TreeNode{View: subView, Parent: cur.Value.(*TreeNode)}
			cur.Value.(*TreeNode).Child = append(cur.Value.(*TreeNode).Child, node)
			stack.PushBack(node)
		}
	}
	return treeRoot
}

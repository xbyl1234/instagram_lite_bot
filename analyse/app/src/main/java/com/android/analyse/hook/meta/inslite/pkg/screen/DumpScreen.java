package com.android.analyse.hook.meta.inslite.pkg.screen;

import com.alibaba.fastjson2.JSONObject;
import com.android.analyse.hook.meta.common.ClassloaderHook;
import com.android.analyse.hook.meta.inslite.ClassNames;
import com.android.analyse.hook.meta.inslite.FieldName;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.common.tools.hooker.HookTools;

import java.lang.reflect.Method;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class DumpScreen {
    public static class TreeNode {
        Object view;
        public TreeNode parent;
        public List<TreeNode> child = new ArrayList<>();
    }

    static Map<String, String> viewClzMap = new HashMap<String, String>() {{
        put(ClassNames.subWrapScreen, "screen0");
        put(ClassNames.subWrapScreen1, "screen1");
        put(ClassNames.subWrapScreen2, "screen2");
        put(ClassNames.subWrapScreen3, "screen3");
        put(ClassNames.subWrapScreen9, "screen9");
        put(ClassNames.subWrapScreen13, "screen13");
        put(ClassNames.subWrapScreen19, "screen19");
    }};

    public static String getScreenLikeResId(Object screen) {
        return "" + HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.BaseScreenBean), screen, FieldName.BaseScreenBeanLikeActionResourceId);
    }

    public static String getWindowName(Object screen) {
        return "" + HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.BaseScreenBean), screen, FieldName.BaseScreenBeanWindowName);
    }

    public static String getScreenName(Object screen) {
        return "" + HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.ScreenDecodeBodyBean), screen, FieldName.ScreenDecodeBodyScreenName);
    }

    public static String getScreenId(Object screen) {
        return "" + HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.BaseScreenBean), screen, FieldName.BaseScreenBeanScreenId);
    }

    static int nameIdx = 0;

    public static String getScreenNodeName(Object root) {
        String name = viewClzMap.get(root.getClass().getName());
        if (name == null) {
            name = root.getClass().getName();
        }
        nameIdx += 1;
        name += "_" + nameIdx;
        name += "_" + root.toString().substring(6);
        name += "_" + getWindowName(root);
        name += "_" + getScreenLikeResId(root);
        return name;
    }

    public static JSONObject tree2Json(TreeNode root) {
        JSONObject json = new JSONObject();
        for (int i = 0; i < root.child.size(); i++) {
            json.put(getScreenNodeName(root.child.get(i).view), tree2Json(root.child.get(i)));
        }
        return json;
    }

    public static String tree2string(TreeNode root) {
        nameIdx = 0;
        if (root == null) {
            return "root is null";
        }
        try {
            JSONObject json = new JSONObject();
            json.put(getScreenNodeName(root.view), tree2Json(root));
            return json.toString();
        } catch (Throwable e) {
            return "error: " + e;
        }
    }

    public static TreeNode dumpViewTree(Object root) throws Throwable {
        Class subWrapScreen2Clz = ClassloaderHook.GetClass(ClassNames.subWrapScreen2);
        Class implClz = ClassloaderHook.GetClass(ClassNames.subWrapScreenImpl);
        Method getAllSubScreen = subWrapScreen2Clz.getMethod(MethodNames.SubWrapScreen2_getAllSubScreen);

        List<TreeNode> stack = new ArrayList<>();
        TreeNode treeRoot = new TreeNode();
        treeRoot.view = root;
        stack.add(treeRoot);
        while (!stack.isEmpty()) {
            TreeNode cur = stack.get(stack.size() - 1);
            stack.remove(stack.size() - 1);
            Object view = cur.view;
            if (implClz.isAssignableFrom(view.getClass())) {
                ArrayList allSub = (ArrayList) getAllSubScreen.invoke(view);
                if (allSub == null) {
                    continue;
                }
                for (int i = 0; i < allSub.size(); i++) {
                    Object subView = allSub.get(i);
                    TreeNode node = new TreeNode();
                    node.view = subView;
                    node.parent = cur;
                    cur.child.add(node);
                    stack.add(node);
                }
            }
        }
        return treeRoot;
    }

}

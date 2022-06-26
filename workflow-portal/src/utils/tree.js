export function buildTree(folders, projectName, projectId) {
  /*
  Folder object contains the following information
    {
    "id": "42d3858b-7954-4a60-825c-87265045e5fd", 
    "name": "test 3", 
    "path": "/test3", 
    "author": "dongoctuan.0101@gmail.com", 
    "created_at": "2021-06-29T06:46:19.526106Z", 
    "updated_at": "2021-06-29T06:46:19.526106Z", 
    "projectId": "4e07845b-5569-4871-b19e-1c8ac90e3052", 
    "project_name": "Bio analysis", 
    "class": "folder"
    }

  List folders in flat structure, so need to build a tree hierarchy structure
  to easy use in ElementPlus el-tree component. The hierarchy is based on path
  property of each folder object.
  */
  if (!folders) {
    return null;
  }

  const tree = [];
  const lookup = {};
  // Add the first root node
  lookup["/"] = {
    id: projectId,
    name: projectName,
    path: "/",
    children: [],
    class: "project",
    project_id: projectId,
    project_name: projectName,
  };

  folders.sort((folder1, folder2) => {
    if (folder1.path === folder2.path) return 0;
    if (folder1.path.length <= folder2.path.length) return -1;
    return 1;
  });

  for (const folder of folders) {
    const path = folder.path;
    let parentPath = path.slice(0, path.lastIndexOf("/"));
    if (parentPath === "") parentPath = "/";

    let parent = lookup["/"]
    if (parent.path !== parentPath) {
      let stack = [];
      stack = stack.concat(parent.children);
      while (stack.length) {
        const node = stack.pop();
        if (node.path === parentPath) {
          parent = node;
          break;
        }
        stack = stack.concat(node.children);
      }
    }

    parent.children.push(
      Object.assign({
        project_id: projectId, project_name: projectName, children: []
      }, folder));
  };

  tree.push(lookup["/"]);
  return tree;
}

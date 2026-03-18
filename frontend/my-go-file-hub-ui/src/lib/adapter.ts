import type { FileNode } from '../types';

const getFileType = (name: string, isDir: boolean): string => {
  if (isDir) return 'folder';
  const ext = name.split('.').pop()?.toLowerCase();
  if (['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext!)) return 'image';
  if (['mp4', 'mkv', 'avi', 'mov'].includes(ext!)) return 'video';
  if (['pdf', 'doc', 'docx', 'txt', 'md'].includes(ext!)) return 'text';
  return 'file';
};

/**
 * 将后端返回的 FileInfo 转换为前端使用的 FileNode
 * 兼容处理：蛇形命名 (is_dir)、大驼峰 (IsDir)、小驼峰 (isDir)
 */
export const adaptFileNode = (
  file: any, 
  currentPath: string
): FileNode => {
  // 兼容处理字段名
  const name = file.name || file.Name || "";
  const isDir = typeof file.is_dir === 'boolean' ? file.is_dir : 
                (typeof file.isDir === 'boolean' ? file.isDir : (file.IsDir || false));
  const size = typeof file.size === 'number' ? file.size : (file.Size || 0);
  const modTime = file.mod_time || file.modTime || file.ModTime || "";

  // 确保路径拼接正确，去除多余的斜杠
  const fullPath = `${currentPath}/${name}`.replace(/\/+/g, '/');
  
  return {
    name,
    isDir,
    size,
    modTime,
    id: fullPath,
    fullPath: fullPath,
    extension: name.includes('.') ? name.split('.').pop()! : '',
    mimeType: getFileType(name, isDir)
  };
};

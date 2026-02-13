import type { FileInfo, FileNode } from '../types';

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
 */
export const adaptFileNode = (
  file: FileInfo, 
  currentPath: string
): FileNode => {
  // 确保路径拼接正确，去除多余的斜杠
  const fullPath = `${currentPath}/${file.name}`.replace(/\/+/g, '/');
  
  return {
    ...file,
    id: fullPath,
    fullPath: fullPath,
    extension: file.name.includes('.') ? file.name.split('.').pop()! : '',
    mimeType: getFileType(file.name, file.isDir)
  };
};

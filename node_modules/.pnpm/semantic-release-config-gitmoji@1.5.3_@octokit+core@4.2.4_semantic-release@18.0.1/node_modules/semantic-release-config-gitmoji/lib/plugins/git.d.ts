import type { PluginSpec } from 'semantic-release';
import type { GitPluginOpts } from '../type';
/**
 * git
 * @param options
 */
declare const git: (options?: GitPluginOpts) => PluginSpec;
export default git;

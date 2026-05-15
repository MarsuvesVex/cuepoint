import type { PluginSpec } from 'semantic-release';
import type { GithubPluginOpts } from '../type';
/**
 * github
 * @param options
 */
declare const github: (options?: GithubPluginOpts) => PluginSpec;
export default github;

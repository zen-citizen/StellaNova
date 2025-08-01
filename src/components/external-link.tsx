import { forwardRef, type AnchorHTMLAttributes } from "react";

type ExternalLinkProps = AnchorHTMLAttributes<HTMLAnchorElement>;

export const ExternalLink = forwardRef<HTMLAnchorElement, ExternalLinkProps>(
  (props, ref) => (
    <a target="__blank" rel="noopener noreferrer" {...props} ref={ref} />
  )
);

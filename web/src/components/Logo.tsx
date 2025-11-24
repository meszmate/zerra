export default function Logo({
    className,
}: {
    className: string,
}) {
    return (
        <svg className={className} viewBox="0 0 305 301" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect y="70" width="231" height="231" rx="28" fill="#2DEA63" fill-opacity="0.5" />
            <rect x="74" width="231" height="231" rx="28" fill="#2DEA63" />
        </svg>
    )
}
